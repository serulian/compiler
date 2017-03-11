// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webidl

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/webidl/parser"

	"github.com/streamrail/concurrent-map"
)

// TypeCollapser represents a collapsed set of the types that will be emitted
// into the type graph from the imported WebIDL modules. This is necessary
// because in order to ensure compatibility in the flat namespace of WebIDL,
// we need to merge types with the same name.
type TypeCollapser struct {
	irg                *WebIRG
	typesEncountered   cmap.ConcurrentMap
	globalDeclarations []IRGDeclaration
}

func createTypeCollapser(irg *WebIRG, modifier compilergraph.GraphLayerModifier) *TypeCollapser {
	tc := &TypeCollapser{
		irg:                irg,
		typesEncountered:   cmap.New(),
		globalDeclarations: make([]IRGDeclaration, 0),
	}

	tc.populate(modifier)
	return tc
}

type ctHandler func(ct *CollapsedType)
type gdHandler func(gd IRGDeclaration)

func (tc *TypeCollapser) ForEachType(handler ctHandler) {
	// TODO: ||
	for ct := range tc.Types() {
		handler(ct)
	}
}

func (tc *TypeCollapser) ForEachGlobalDeclaration(handler gdHandler) {
	// TODO: ||
	for _, globalDecl := range tc.globalDeclarations {
		handler(globalDecl)
	}
}

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
				return &CollapsedType{
					RootNode:     globalDeclaration,
					Declarations: make([]IRGDeclaration, 0, 1),
					Name:         name,
					Serializable: false,
					ParentTypes:  map[string]IRGDeclaration{},

					Operators:       map[string]IRGAnnotation{},
					Members:         map[string]IRGMember{},
					Specializations: map[MemberSpecialization]IRGMember{},
				}
			}

			return valueInMap
		})

	ct, _ := tc.GetType(name)
	return ct
}

func (tc *TypeCollapser) populate(modifier compilergraph.GraphLayerModifier) {
	collapseDeclaration := func(key interface{}, value interface{}) bool {
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

type CollapsedType struct {
	RootNode     compilergraph.GraphNode
	Declarations []IRGDeclaration

	Name                   string
	ConstructorAnnotations []IRGAnnotation

	Serializable bool
	ParentTypes  map[string]IRGDeclaration

	Operators       map[string]IRGAnnotation
	Members         map[string]IRGMember
	Specializations map[MemberSpecialization]IRGMember
}

func (ct *CollapsedType) RegisterOperator(name string, opAnnotation IRGAnnotation) bool {
	if _, exists := ct.Operators[name]; exists {
		return false
	}

	ct.Operators[name] = opAnnotation
	return true
}

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
