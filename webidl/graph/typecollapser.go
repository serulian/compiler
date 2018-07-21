// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/webidl/parser"

	cmap "github.com/streamrail/concurrent-map"
)

// TypeCollapser represents a collapsed set of the types that will be emitted
// into the type graph from the imported WebIDL modules. This is necessary
// because in order to ensure compatibility in the flat namespace of WebIDL,
// we need to merge types with the same name.
type TypeCollapser struct {
	irg                    *WebIRG
	typesEncountered       cmap.ConcurrentMap
	collapsedTypesByNodeID cmap.ConcurrentMap
	globalDeclarations     []IRGDeclaration
}

// createTypeCollapser returns a new populated type collapser. Note that this call modifies the
// graph, and therefore should only be invoked once.
func createTypeCollapser(irg *WebIRG, modifier compilergraph.GraphLayerModifier) *TypeCollapser {
	tc := emptyTypeCollapser(irg)
	tc.populate(modifier)
	return tc
}

func emptyTypeCollapser(irg *WebIRG) *TypeCollapser {
	tc := &TypeCollapser{
		irg:                    irg,
		typesEncountered:       cmap.New(),
		collapsedTypesByNodeID: cmap.New(),
		globalDeclarations:     make([]IRGDeclaration, 0),
	}
	return tc
}

type ctHandler func(ct *CollapsedType)
type gdHandler func(gd IRGDeclaration)

// ForEachType invokes the given handler for each collapsed type. Note that the handlers will
// be invoked in parallel via goroutines, so care must be taken if access any shared resources.
func (tc *TypeCollapser) ForEachType(handler ctHandler) {
	process := func(key interface{}, value interface{}, cancel compilerutil.CancelFunction) bool {
		handler(key.(*CollapsedType))
		return true
	}

	workqueue := compilerutil.Queue()
	for ct := range tc.Types() {
		workqueue.Enqueue(ct, ct, process)
	}
	workqueue.Run()
}

// ForEachGlobalDeclaration invokes the given handler for each global decl. Note that the handlers will
// be invoked in parallel via goroutines, so care must be taken if access any shared resources.
func (tc *TypeCollapser) ForEachGlobalDeclaration(handler gdHandler) {
	process := func(key interface{}, value interface{}, cancel compilerutil.CancelFunction) bool {
		handler(key.(IRGDeclaration))
		return true
	}

	workqueue := compilerutil.Queue()
	for _, globalDecl := range tc.globalDeclarations {
		workqueue.Enqueue(globalDecl, globalDecl, process)
	}
	workqueue.Run()
}

// Types returns all collapsed types found in the WebIDL IRG.
func (tc *TypeCollapser) Types() chan *CollapsedType {
	ch := make(chan *CollapsedType, 10)
	go (func() {
		for tuple := range tc.typesEncountered.IterBuffered() {
			ch <- tuple.Val.(*CollapsedType)
		}
		close(ch)
	})()
	return ch
}

// GetTypeForNodeID returns the collapsed type with the given matching Node ID.
func (tc *TypeCollapser) GetTypeForNodeID(nodeID compilergraph.GraphNodeId) (*CollapsedType, bool) {
	collapsed, found := tc.collapsedTypesByNodeID.Get(string(nodeID))
	if !found {
		return nil, false
	}

	return collapsed.(*CollapsedType), true
}

// GetType returns the collapsed type matching the given name, if any.
func (tc *TypeCollapser) GetType(name string) (*CollapsedType, bool) {
	ct, found := tc.typesEncountered.Get(name)
	if !found {
		return &CollapsedType{}, false
	}

	return ct.(*CollapsedType), true
}

func (tc *TypeCollapser) getType(modifier compilergraph.GraphLayerModifier, name string) *CollapsedType {
	// Perform an upsert into the concurrent map, only creating a new global declaration node and
	// associated CollapsedType struct on the first existance of the name encountered.
	tc.typesEncountered.Upsert(name, nil,
		func(exist bool, valueInMap interface{}, newValue interface{}) interface{} {
			if !exist {
				globalDeclaration := modifier.CreateNode(parser.NodeTypeGlobalDeclaration).AsNode()
				collapsedType := &CollapsedType{
					RootNode:     globalDeclaration,
					Declarations: make([]IRGDeclaration, 0, 1),
					Name:         name,
					Serializable: false,
					ParentTypes:  map[string]IRGDeclaration{},

					Operators:       map[string]IRGAnnotation{},
					Members:         map[string]IRGMember{},
					Specializations: map[MemberSpecialization]IRGMember{},
				}

				tc.collapsedTypesByNodeID.Set(string(globalDeclaration.GetNodeId()), collapsedType)
				return collapsedType
			}

			return valueInMap
		})

	ct, _ := tc.GetType(name)
	return ct
}

// populate populates the type collapser (and IRG) with a root type node for each shared
// name defined via the various declarations in the IRG. For example, two declarations
// named `Number` will result in a single *CollapsedType with name `Number`, a backing
// node in the IRG, and the list of those declarations contained within.
func (tc *TypeCollapser) populate(modifier compilergraph.GraphLayerModifier) {
	collapseDeclaration := func(key interface{}, value interface{}, cancel compilerutil.CancelFunction) bool {
		name := key.(string)
		declaration := value.(IRGDeclaration)

		// Retrieve the shared collapsed type for this type name.
		ct := tc.getType(modifier, name)

		// Add the declaration to the declaration list, along with its parent type
		// and whether it is marked serializable.
		ct.Declarations = append(ct.Declarations, declaration)

		if parentType, hasParentType := declaration.ParentType(); hasParentType {
			ct.ParentTypes[parentType] = declaration
		}

		if !ct.Serializable && declaration.IsSerializable() {
			ct.Serializable = true
		}

		// Add the constructor annotation(s), if any.
		for _, annotation := range declaration.GetAnnotations(CONSTRUCTOR_ANNOTATION) {
			ct.ConstructorAnnotations = append(ct.ConstructorAnnotations, annotation)
		}

		return true
	}

	// Enqueue each declaration to be collapsed, with the key being the declaration name, to ensure
	// that we never concurrently work on a type with the same name.
	workqueue := compilerutil.Queue()
	for _, module := range tc.irg.GetModules() {
		for _, declaration := range module.Declarations() {
			// If the type is marked as [Global] then it defines an "interface" whose members
			// get added to the global context, and, therefore, not a real type.
			if declaration.HasOneAnnotation(GLOBAL_CONTEXT_ANNOTATIONS...) {
				tc.globalDeclarations = append(tc.globalDeclarations, declaration)
			} else {
				workqueue.Enqueue(declaration.Name(), declaration, collapseDeclaration)
			}
		}
	}
	workqueue.Run()
}

// CollapsedType represents a single named type in the WebIDL that has been collapsed
// from (possibly multiple) declarations.
type CollapsedType struct {
	RootNode     compilergraph.GraphNode // The root node for this type in the IRG.
	Declarations []IRGDeclaration        // The declarations that created this type.

	Name                   string          // The name of this type. Unique amongst all CollapsedType's.
	ConstructorAnnotations []IRGAnnotation // The constructor(s) for this type, if any.

	Serializable bool                      // Whether this type is marked as serializable.
	ParentTypes  map[string]IRGDeclaration // The parent type(s) for this type, if any.

	Operators       map[string]IRGAnnotation           // The registered operators.
	Members         map[string]IRGMember               // The registered members.
	Specializations map[MemberSpecialization]IRGMember // The registered specializations.
}

// RegisterOperator registers an operator with the given name and annotation, returning
// true if this is the first occurance of the operator under the collapsed type.
func (ct *CollapsedType) RegisterOperator(name string, opAnnotation IRGAnnotation) bool {
	if _, exists := ct.Operators[name]; exists {
		return false
	}

	ct.Operators[name] = opAnnotation
	return true
}

// RegisterMember registers a member, returning true if this is the first occurance of the
// member under the collapsed type. If another member of the same name exists *and* its
// signature does not match, an error is reported on the reporter.
func (ct *CollapsedType) RegisterMember(member IRGMember, reporter typegraph.IssueReporter) bool {
	name, _ := member.Name()
	if existingMember, exists := ct.Members[name]; exists {
		if existingMember.Signature() != member.Signature() {
			reporter.ReportError(member.GraphNode, "Member '%s' redefined under type '%s' but with a different signature", name, ct.Name)
		}

		return false
	}

	ct.Members[name] = member
	return true
}

// RegisterSpecialization registers a specialization member, returning true if this is the first occurance of the
// member under the collapsed type. If another member of the same specialization exists *and* its
// signature does not match, an error is reported on the reporter.
func (ct *CollapsedType) RegisterSpecialization(member IRGMember, reporter typegraph.IssueReporter) bool {
	specialization, _ := member.Specialization()
	if existingMember, exists := ct.Specializations[specialization]; exists {
		if existingMember.Signature() != member.Signature() {
			reporter.ReportError(member.GraphNode, "'%s' redefined under type '%s' but with a different signature", specialization, ct.Name)
		}

		return false
	}

	ct.Specializations[specialization] = member
	return true
}

// SourceRanges returns the ranges for this collapsed type.
func (ct *CollapsedType) SourceRanges() []compilercommon.SourceRange {
	var ranges = make([]compilercommon.SourceRange, 0, len(ct.Declarations))
	for _, declaration := range ct.Declarations {
		sourceRange, hasSourceRange := declaration.SourceRange()
		if hasSourceRange {
			ranges = append(ranges, sourceRange)
		}
	}
	return ranges
}
