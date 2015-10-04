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

	// Add members (along full inheritance)
	for _, srgType := range g.GetTypes() {
		// TODO: add inheritance and cycle checking
		for _, member := range srgType.Members() {
			typeDecl := TGTypeDecl{t.getTypeNodeForSRGType(srgType), t}
			success := t.buildMemberNode(typeDecl, member)
			if !success {
				result.Status = false
			}
		}
	}

	// Constraint check all type references.
	// TODO: this.

	// If the result is not true, collect all the errors found.
	if !result.Status {
		it := t.layer.StartQuery().
			With(NodePredicateError).
			BuildNodeIterator(NodePredicateSource)

		for it.Next() {
			node := it.Node()

			// Lookup the location of the SRG source node.
			srgSourceNode := g.GetNode(compilergraph.GraphNodeId(it.Values()[NodePredicateSource]))
			location := g.NodeLocation(srgSourceNode)

			// Add the error.
			errNode := node.GetNode(NodePredicateError)
			msg := errNode.Get(NodePredicateErrorMessage)
			result.Errors = append(result.Errors, compilercommon.NewSourceError(location, msg))
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

type genericWork struct {
	generic    srg.SRGGeneric
	parentType srg.SRGType
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
		_, ok := t.buildGenericNode(data.generic, parentType, NodePredicateTypeGeneric)
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
		for _, srgGeneric := range srgType.Generics() {
			genericKey := workKey{srgType.Node(), srgGeneric.Name(), "Generic"}
			workqueue.Enqueue(genericKey, genericWork{srgGeneric, srgType}, buildTypeGeneric, typeKey)
		}
	}

	// Run the queue to construct the full set of types and generics.
	return workqueue.Run()
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

	return workqueue.Run()
}

// decorateWithError decorates the given node with an associated error node.
func (t *TypeGraph) decorateWithError(node compilergraph.GraphNode, message string, args ...interface{}) {
	errorNode := t.layer.CreateNode(NodeTypeError)
	errorNode.Decorate(NodePredicateErrorMessage, fmt.Sprintf(message, args...))
	node.Connect(NodePredicateError, errorNode)
}
