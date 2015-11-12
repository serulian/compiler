// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"strconv"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
)

type genericKind int

const (
	typeDeclGeneric genericKind = iota
	typeMemberGeneric
)

// resolveGenericConstraint decorates a generic node with its defined constraint. Returns false if the constraint
// is invalid due to a resolution error.
func (t *TypeGraph) resolveGenericConstraint(generic srg.SRGGeneric, genericNode compilergraph.GraphNode) bool {
	constraintType, valid := t.resolvePossibleType(genericNode, generic.GetConstraint)
	genericNode.DecorateWithTagged(NodePredicateGenericSubtype, constraintType)
	return valid
}

// buildGenericNode adds a new generic node to the specified type or type membe node for the given SRG generic.
func (t *TypeGraph) buildGenericNode(generic srg.SRGGeneric, index int, kind genericKind, parentNode compilergraph.GraphNode, parentPredicate string) (compilergraph.GraphNode, bool) {
	// Ensure that there exists no other generic with this name under the parent node.
	_, exists := parentNode.StartQuery().
		Out(parentPredicate).
		Has(NodePredicateGenericName, generic.Name()).
		TryGetNode()

	// Create the generic node.
	genericNode := t.layer.CreateNode(NodeTypeGeneric)
	genericNode.Decorate(NodePredicateGenericName, generic.Name())
	genericNode.Decorate(NodePredicateGenericIndex, strconv.Itoa(index))
	genericNode.Decorate(NodePredicateGenericKind, strconv.Itoa(int(kind)))

	genericNode.Connect(NodePredicateSource, generic.Node())

	// Add the generic to the parent node.
	parentNode.Connect(parentPredicate, genericNode)

	// Mark the generic with an error if it is repeated.
	if exists {
		t.decorateWithError(genericNode, "Generic '%s' is already defined", generic.Name())
		return genericNode, false
	}

	return genericNode, true
}
