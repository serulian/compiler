// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
)

// scopeContext represents the currently operating context for scoping, allowing for
// scope-specific overrides of such items as types of expressions.
type scopeContext struct {
	// rootNode holds the root node under which we are scoping.
	rootNode compilergraph.GraphNode

	// accessOption is the access option for the current scope.
	accessOption scopeAccessOption

	// staticDependencyCollector defines a helper for collecting all members accessed or called
	// statically the current scope.
	staticDependencyCollector *staticDependencyCollector

	// dynamicDependencyCollector defines a helper for collecting all names accessed dynamically in
	// the current scope.
	dynamicDependencyCollector *dynamicDependencyCollector

	// overrideTypes is (if not nil) the map of the overridden type for an expression
	// under this context.
	overrideTypes *map[compilergraph.GraphNodeId]typegraph.TypeReference

	// parentImplemented holds the immediate parent implemented node.
	// The parent implemented node is the parent type member or function under which
	// we are scoping the implementation.
	parentImplemented compilergraph.GraphNode

	// parentBreakable holds a reference to the parent node to which a `break` statement
	// can jump, if any.
	parentBreakable *compilergraph.GraphNode

	// parentContinuable holds a reference to the parent node to which a `continue` statement
	// can jump, if any.
	parentContinuable *compilergraph.GraphNode

	// rootLabelSet is the set of extra labels for the root node.
	rootLabelSet *statementLabelSet
}

// getParentContainer returns the parent type member, module member or property getter/setter
// under which we are scoping, if any.
func (sc scopeContext) getParentContainer(g *srg.SRG) (srg.SRGImplementable, bool) {
	return g.AsImplementable(sc.rootNode)
}

// getTypeOverride returns the type override for the given expression node, if any.
func (sc scopeContext) getTypeOverride(exprNode compilergraph.GraphNode) (typegraph.TypeReference, bool) {
	if sc.overrideTypes == nil {
		return typegraph.TypeReference{}, false
	}

	ot := *sc.overrideTypes
	value, found := ot[exprNode.NodeId]
	return value, found
}

// withContinuable returns the scope context with the parent continuable and breakable nodes
// set to that given.
func (sc scopeContext) withContinuable(node compilergraph.GraphNode) scopeContext {
	return scopeContext{
		rootNode:                   sc.rootNode,
		staticDependencyCollector:  sc.staticDependencyCollector,
		dynamicDependencyCollector: sc.dynamicDependencyCollector,

		accessOption:      sc.accessOption,
		overrideTypes:     sc.overrideTypes,
		parentImplemented: sc.parentImplemented,

		rootLabelSet: sc.rootLabelSet,

		parentBreakable:   &node,
		parentContinuable: &node,
	}
}

// withBreakable returns the scope context with the parent breakable node set to that
// given.
func (sc scopeContext) withBreakable(node compilergraph.GraphNode) scopeContext {
	return scopeContext{
		rootNode:                   sc.rootNode,
		staticDependencyCollector:  sc.staticDependencyCollector,
		dynamicDependencyCollector: sc.dynamicDependencyCollector,

		accessOption:      sc.accessOption,
		overrideTypes:     sc.overrideTypes,
		parentImplemented: sc.parentImplemented,
		parentContinuable: sc.parentContinuable,

		rootLabelSet: sc.rootLabelSet,

		parentBreakable: &node,
	}
}

// withImplemented returns the scope context with the parent implemented node set to that
// given.
func (sc scopeContext) withImplemented(node compilergraph.GraphNode) scopeContext {
	return scopeContext{
		rootNode:                   sc.rootNode,
		staticDependencyCollector:  sc.staticDependencyCollector,
		dynamicDependencyCollector: sc.dynamicDependencyCollector,

		accessOption:      sc.accessOption,
		overrideTypes:     sc.overrideTypes,
		parentBreakable:   sc.parentBreakable,
		parentContinuable: sc.parentContinuable,

		rootLabelSet: sc.rootLabelSet,

		parentImplemented: node,
	}
}

// withAccess returns the scope context with the access option set to that given.
func (sc scopeContext) withAccess(access scopeAccessOption) scopeContext {
	return scopeContext{
		rootNode:                   sc.rootNode,
		staticDependencyCollector:  sc.staticDependencyCollector,
		dynamicDependencyCollector: sc.dynamicDependencyCollector,

		overrideTypes:     sc.overrideTypes,
		parentImplemented: sc.parentImplemented,
		parentBreakable:   sc.parentBreakable,
		parentContinuable: sc.parentContinuable,

		rootLabelSet: sc.rootLabelSet,

		accessOption: access,
	}
}

// withTypeOverride returns the scope context with the type of the given expression node
// overridden.
func (sc scopeContext) withTypeOverride(exprNode compilergraph.GraphNode, typeref typegraph.TypeReference) scopeContext {
	overrideTypes := map[compilergraph.GraphNodeId]typegraph.TypeReference{}

	if sc.overrideTypes != nil {
		existing := *sc.overrideTypes
		for key, value := range existing {
			overrideTypes[key] = value
		}
	}

	overrideTypes[exprNode.NodeId] = typeref

	return scopeContext{
		rootNode:                   sc.rootNode,
		staticDependencyCollector:  sc.staticDependencyCollector,
		dynamicDependencyCollector: sc.dynamicDependencyCollector,

		accessOption:      sc.accessOption,
		parentImplemented: sc.parentImplemented,
		parentBreakable:   sc.parentBreakable,
		parentContinuable: sc.parentContinuable,

		rootLabelSet: sc.rootLabelSet,

		overrideTypes: &overrideTypes,
	}
}
