// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/srg"
)

// Build builds the type graph from the SRG used to initialize it.
func (t *TypeGraph) build(g *srg.SRG) *Result {
	result := &Result{
		Status:   true,
		Warnings: make([]*compilercommon.SourceWarning, 0),
		Errors:   make([]*compilercommon.SourceError, 0),
		Graph:    t,
	}

	// Translate modules.
	t.translateModules()

	// Translate types and generics.
	if !t.translateTypesAndGenerics() {
		result.Status = false
	}

	// Resolve generic constraints. We do this outside of adding generics as they may depend
	// on other generics.
	if !t.resolveGenericConstraints() {
		result.Status = false
	}

	// Load the operators map. Requires the types loaded as it performs lookups of certain types (int, etc).
	t.buildOperatorDefinitions()

	// Translate module members.
	t.translateModuleMembers()

	// Add type members (along full inheritance)
	if !t.translateTypeMembers() {
		result.Status = false
	}

	// Constraint check all type references.
	if result.Status && !t.verifyTypeReferences() {
		result.Status = false
	}

	// If the result is not true, collect all the errors found.
	if !result.Status {
		it := t.layer.StartQuery().
			With(NodePredicateError).
			BuildNodeIterator(NodePredicateSource)

		for it.Next() {
			node := it.Node()

			// Lookup the location of the SRG source node.
			if sourceNodeId, ok := it.Values()[NodePredicateSource]; ok {
				srgSourceNode := g.GetNode(compilergraph.GraphNodeId(sourceNodeId))
				location := g.NodeLocation(srgSourceNode)

				// Add the error.
				errNode := node.GetNode(NodePredicateError)
				msg := errNode.Get(NodePredicateErrorMessage)
				result.Errors = append(result.Errors, compilercommon.NewSourceError(location, msg))
			} else {
				panic(fmt.Sprintf("Error on non-sourced node: %v", it.Node()))
			}
		}
	}

	return result
}

// workKey is a struct representing a key for a concurrent work job during construction.
type workKey struct {
	ParentNode compilergraph.GraphNode
	Name       string
	Kind       string
}

// genericWork holds data for a generic translation.
type genericWork struct {
	generic    srg.SRGGeneric
	index      int
	parentType srg.SRGType
}

// typeMemberWork holds data for type member translations.
type typeMemberWork struct {
	srgType          srg.SRGType
	parentReferences []TypeReference
}

// moduleMemberWork holds data for a module member translations.
type moduleMemberWork struct {
	srgModule srg.SRGModule
	srgMember srg.SRGMember
}

// verifyTypeReferences verifies and validates all type references in the *SRG*.
func (t *TypeGraph) verifyTypeReferences() bool {
	validateTyperef := func(key interface{}, value interface{}) bool {
		srgTypeRef := value.(srg.SRGTypeRef)
		typeref, err := t.BuildTypeRef(srgTypeRef)
		if err != nil {
			issueNode := t.layer.CreateNode(NodeTypeReferenceIssue)
			issueNode.Connect(NodePredicateSource, srgTypeRef.GraphNode)
			t.decorateWithError(issueNode, "%v", err)
			return false
		}

		verr := typeref.Verify()
		if verr != nil {
			issueNode := t.layer.CreateNode(NodeTypeReferenceIssue)
			issueNode.Connect(NodePredicateSource, srgTypeRef.GraphNode)
			t.decorateWithError(issueNode, "%v", verr)
			return false
		}

		return true
	}

	workqueue := compilerutil.Queue()

	for _, srgTypeRef := range t.srg.GetTypeReferences() {
		refKey := workKey{srgTypeRef.GraphNode, "-", "Typeref"}
		workqueue.Enqueue(refKey, srgTypeRef, validateTyperef)
	}

	return workqueue.Run().Status
}

// translateModules translate the modules from the SRG.
func (t *TypeGraph) translateModules() bool {
	buildModule := func(key interface{}, value interface{}) bool {
		return t.buildModuleNode(value.(srg.SRGModule))
	}

	// Enqueue the set of SRG modules.
	workqueue := compilerutil.Queue()

	for _, module := range t.srg.GetModules() {
		moduleKey := workKey{module.GraphNode, module.Name(), "Module"}
		workqueue.Enqueue(moduleKey, module, buildModule)
	}

	return workqueue.Run().Status
}

// translateModuleMembers translate direcy module members from the SRG.
func (t *TypeGraph) translateModuleMembers() bool {
	buildModuleMember := func(key interface{}, value interface{}) bool {
		data := value.(moduleMemberWork)
		module, found := t.getModuleForSRGModule(data.srgModule)
		if !found {
			panic(fmt.Sprintf("Could not find matching TG module for %v", data.srgModule.Name()))
		}

		return t.buildMemberNode(module, data.srgMember)
	}

	// Enqueue the set of SRG module members.
	workqueue := compilerutil.Queue()

	for _, module := range t.srg.GetModules() {
		// Check for direct members under the module.
		moduleMembers := module.GetMembers()
		if len(moduleMembers) == 0 {
			continue
		}

		// Enqueue the members of the module.
		for _, member := range moduleMembers {
			memberKey := workKey{member.GraphNode, member.Name(), "ModuleMember"}
			workqueue.Enqueue(memberKey, moduleMemberWork{module, member}, buildModuleMember)
		}
	}

	return workqueue.Run().Status
}

// translateTypeMembers translates the type members from the SRG, including handling class
// composition/inheritance.
func (t *TypeGraph) translateTypeMembers() bool {
	buildTypeMembers := func(key interface{}, value interface{}) bool {
		data := value.(typeMemberWork)
		typeDecl := TGTypeDecl{t.getTypeNodeForSRGType(data.srgType), t}
		return t.buildMembership(typeDecl, data.srgType, data.parentReferences)
	}

	// Enqueue the full set of SRG types with dependencies on any parent types.
	workqueue := compilerutil.Queue()

	var success = true
	for _, srgType := range t.srg.GetTypes() {
		// The key for inheritance is the node itself.
		typeNode := t.getTypeNodeForSRGType(srgType)
		typeKey := workKey{typeNode, "-", "Members"}

		// Build the list of dependency keys. A job here depends on the members jobs for
		// all types it inherits/composes. This is done to ensure that the entire membership
		// graph for the parent type is in place before we process this type.
		var dependencies = make([]interface{}, 0)
		var parentTypes = make([]TypeReference, 0)

		for _, inheritsRef := range srgType.Inheritance() {
			// Resolve the type to which the inherits points.
			resolvedType, err := t.BuildTypeRef(inheritsRef)
			if err != nil {
				t.decorateWithError(typeNode, "%s", err.Error())
				success = false
			}

			// Ensure the referred type is a class.
			referredType := resolvedType.referredTypeNode()
			if referredType.Kind != NodeTypeClass {
				switch referredType.Kind {
				case NodeTypeGeneric:
					name := referredType.Get(NodePredicateGenericName)
					t.decorateWithError(typeNode, "Type '%s' cannot derive from a generic ('%s')", srgType.Name(), name)

				case NodeTypeInterface:
					name := referredType.Get(NodePredicateTypeName)
					t.decorateWithError(typeNode, "Type '%s' cannot derive from an interface ('%s')", srgType.Name(), name)
				}

				success = false
				continue
			}

			dependencies = append(dependencies, workKey{referredType, "-", "Members"})
			parentTypes = append(parentTypes, resolvedType)
		}

		workqueue.Enqueue(typeKey, typeMemberWork{srgType, parentTypes}, buildTypeMembers, dependencies...)
	}

	// Run the queue to construct the full inheritance.
	result := workqueue.Run()
	if !result.HasCycle {
		return result.Status && success
	} else {
		// TODO(jschorr): If there are two cycles, this will conflate them. We should do actual
		// checking here.
		var types = make([]string, len(result.Cycle))
		for index, key := range result.Cycle {
			decl := TGTypeDecl{key.(workKey).ParentNode, t}
			types[index] = decl.Name()
		}

		t.decorateWithError(result.Cycle[0].(workKey).ParentNode, "A cycle was detected in the inheritance of types: %v", types)
		return false
	}
}

// translateTypesAndGenerics translates the types and their generics defined in the SRG into
// the type graph.
func (t *TypeGraph) translateTypesAndGenerics() bool {
	buildTypeNode := func(key interface{}, value interface{}) bool {
		return t.buildTypeNode(value.(srg.SRGType))
	}

	buildTypeGeneric := func(key interface{}, value interface{}) bool {
		data := value.(genericWork)
		parentType := t.getTypeNodeForSRGType(data.parentType)
		_, ok := t.buildGenericNode(data.generic, data.index, typeDeclGeneric, parentType, NodePredicateTypeGeneric)
		return ok
	}

	// Enqueue the full set of SRG types and generics to be translated into the type graph.
	workqueue := compilerutil.Queue()

	for _, srgType := range t.srg.GetTypes() {
		// Enqueue the SRG type. Key is the name of the type under its module, as we want to
		// make sure we can check for duplicate names.
		typeKey := workKey{srgType.Module().Node(), srgType.Name(), "Type"}
		workqueue.Enqueue(typeKey, srgType, buildTypeNode)

		// Enqueue the type's generics. The key is the generic name under the type, and the
		// dependency is that the parent type has been constructed.
		for index, srgGeneric := range srgType.Generics() {
			genericKey := workKey{srgType.Node(), srgGeneric.Name(), "Generic"}
			workqueue.Enqueue(genericKey, genericWork{srgGeneric, index, srgType}, buildTypeGeneric, typeKey)
		}
	}

	// Run the queue to construct the full set of types and generics.
	return workqueue.Run().Status
}

// resolveGenericConstraints resolves all constraints defined on generics in the type graph.
func (t *TypeGraph) resolveGenericConstraints() bool {
	resolveGenericConstraint := func(key interface{}, value interface{}) bool {
		srgGeneric := value.(srg.SRGGeneric)
		genericNode := t.getGenericNodeForSRGGeneric(srgGeneric)
		return t.resolveGenericConstraint(srgGeneric, genericNode)
	}

	// Enqueue all generics in the SRG under types. Those under type members will be handled
	// during type member construction.
	workqueue := compilerutil.Queue()

	for _, srgGeneric := range t.srg.GetTypeGenerics() {
		workqueue.Enqueue(srgGeneric, srgGeneric, resolveGenericConstraint)
	}

	return workqueue.Run().Status
}

// decorateWithError decorates the given node with an associated error node.
func (t *TypeGraph) decorateWithError(node compilergraph.GraphNode, message string, args ...interface{}) {
	errorNode := t.layer.CreateNode(NodeTypeError)
	errorNode.Decorate(NodePredicateErrorMessage, fmt.Sprintf(message, args...))
	node.Connect(NodePredicateError, errorNode)
}
