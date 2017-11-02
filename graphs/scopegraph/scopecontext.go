// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
)

// scopeAccessOption defines the kind of access under which the scope
// exists.
type scopeAccessOption int

const (
	scopeGetAccess scopeAccessOption = iota
	scopeSetAccess
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
	staticDependencyCollector staticDependencyCollector

	// dynamicDependencyCollector defines a helper for collecting all names accessed dynamically in
	// the current scope.
	dynamicDependencyCollector dynamicDependencyCollector

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

	// allowAgentConstructions is (if not nil) the "set" of agent types that can be constructed
	// under the current context. This prevents agents from being constructed in code where they
	// will not be *immediately* given to the constructor of their composing type, thus ensuring
	// the `principal` back-reference is available immediately and therefore not breaking type
	// safety.
	allowAgentConstructions *map[typegraph.TypeReference]bool

	// localScopeNamesCache is (if not nil) an immutable map of the cached names found in the local
	// scope. Note that as this is a cache, it is not *guarenteed* to have all local names. It merely
	// makes looking up of local names that we know to add, much faster.
	localScopeNamesCache compilerutil.ImmutableMap
}

// lookupLocalScopeName checks the local name *cache* for the given name returning it if found.
// If the name is *not* found, callers should still execute the slower path via the SRG, as
// the cache is not guarenteed to contain all the names. The use of this cache results in
// significant speedups in scoping, as, at the time it was written, name lookup in the graph
// is incredibly slow in comparison to a map lookup.
func (sc scopeContext) lookupLocalScopeName(name string) (namedScopeInfo, bool) {
	localScopeNamesCache := sc.localScopeNamesCache
	if localScopeNamesCache == nil {
		return namedScopeInfo{}, false
	}

	scope, ok := localScopeNamesCache.Get(name)
	if !ok {
		return namedScopeInfo{}, false
	}

	return scope.(namedScopeInfo), true
}

// withLocalNamed adds the given named node to the local named cache.
func (sc scopeContext) withLocalNamed(node compilergraph.GraphNode, builder *scopeBuilder) scopeContext {
	localScopeNamesCache := sc.localScopeNamesCache
	if localScopeNamesCache == nil {
		localScopeNamesCache = compilerutil.NewImmutableMap()
	}

	scope, err := builder.processSRGNameOrInfo(builder.sg.srg.ScopeNameForNode(node))
	if err != nil {
		return sc
	}

	name, hasName := scope.Name()
	if !hasName {
		return sc
	}

	return scopeContext{
		rootNode:                   sc.rootNode,
		staticDependencyCollector:  sc.staticDependencyCollector,
		dynamicDependencyCollector: sc.dynamicDependencyCollector,

		accessOption:            sc.accessOption,
		overrideTypes:           sc.overrideTypes,
		allowAgentConstructions: sc.allowAgentConstructions,
		parentImplemented:       sc.parentImplemented,
		parentContinuable:       sc.parentContinuable,
		parentBreakable:         sc.parentBreakable,
		rootLabelSet:            sc.rootLabelSet,

		localScopeNamesCache: localScopeNamesCache.Set(name, scope),
	}
}

// getParentContainer returns the parent type member, module member or property getter/setter
// under which we are scoping, if any.
func (sc scopeContext) getParentContainer(g *srg.SRG) (srg.SRGImplementable, bool) {
	return g.AsImplementable(sc.rootNode)
}

// getParentType returns the parent type (and type member) under which we are scoping, if any.
func (sc scopeContext) getParentTypeAndMember(srg *srg.SRG, tdg *typegraph.TypeGraph) (typegraph.TGTypeDecl, typegraph.TGMember, bool, bool) {
	srgImpl, found := sc.getParentContainer(srg)
	if !found {
		return typegraph.TGTypeDecl{}, typegraph.TGMember{}, false, false
	}

	tgMember, tgFound := tdg.GetMemberForSourceNode(srgImpl.ContainingMember().GraphNode)
	if !tgFound {
		return typegraph.TGTypeDecl{}, tgMember, false, false
	}

	tgType, typeFound := tgMember.ParentType()
	return tgType, tgMember, typeFound, true
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

// allowsAgentConstruction returns whether construction of the agent with the given type is allowed
// under the current context.
func (sc scopeContext) allowsAgentConstruction(agentType typegraph.TypeReference) bool {
	if sc.allowAgentConstructions == nil {
		return false
	}

	checkMap := *sc.allowAgentConstructions
	_, exists := checkMap[agentType]
	return exists
}

// withContinuable returns the scope context with the parent continuable and breakable nodes
// set to that given.
func (sc scopeContext) withContinuable(node compilergraph.GraphNode) scopeContext {
	return scopeContext{
		rootNode:                   sc.rootNode,
		staticDependencyCollector:  sc.staticDependencyCollector,
		dynamicDependencyCollector: sc.dynamicDependencyCollector,

		accessOption:            sc.accessOption,
		overrideTypes:           sc.overrideTypes,
		allowAgentConstructions: sc.allowAgentConstructions,
		parentImplemented:       sc.parentImplemented,

		rootLabelSet:         sc.rootLabelSet,
		localScopeNamesCache: sc.localScopeNamesCache,

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

		accessOption:            sc.accessOption,
		overrideTypes:           sc.overrideTypes,
		allowAgentConstructions: sc.allowAgentConstructions,
		parentImplemented:       sc.parentImplemented,
		parentContinuable:       sc.parentContinuable,

		rootLabelSet:         sc.rootLabelSet,
		localScopeNamesCache: sc.localScopeNamesCache,

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

		accessOption:            sc.accessOption,
		overrideTypes:           sc.overrideTypes,
		allowAgentConstructions: sc.allowAgentConstructions,
		parentBreakable:         sc.parentBreakable,
		parentContinuable:       sc.parentContinuable,

		rootLabelSet:         sc.rootLabelSet,
		localScopeNamesCache: sc.localScopeNamesCache,

		parentImplemented: node,
	}
}

// withAccess returns the scope context with the access option set to that given.
func (sc scopeContext) withAccess(access scopeAccessOption) scopeContext {
	return scopeContext{
		rootNode:                   sc.rootNode,
		staticDependencyCollector:  sc.staticDependencyCollector,
		dynamicDependencyCollector: sc.dynamicDependencyCollector,

		overrideTypes:           sc.overrideTypes,
		allowAgentConstructions: sc.allowAgentConstructions,
		parentImplemented:       sc.parentImplemented,
		parentBreakable:         sc.parentBreakable,
		parentContinuable:       sc.parentContinuable,

		rootLabelSet:         sc.rootLabelSet,
		localScopeNamesCache: sc.localScopeNamesCache,

		accessOption: access,
	}
}

// withAllowedAgentConstructionsOf returns the scope context with the agent types composed into
// the given type registered as allowed construction under this context.
func (sc scopeContext) withAllowedAgentConstructionsOf(parentType typegraph.TypeReference) scopeContext {
	if !parentType.IsRefToClass() && !parentType.IsRefToAgent() {
		return sc
	}

	var current = sc
	for _, agent := range parentType.ReferredType().ComposedAgents() {
		current = current.withAllowedAgentConstruction(agent.AgentType())
	}
	return current
}

// withAllowedAgentConstruction returns the scope context with the given agent type added as an allowed
// construction of an agent under this context.
func (sc scopeContext) withAllowedAgentConstruction(agentType typegraph.TypeReference) scopeContext {
	if !agentType.IsRefToAgent() {
		panic("Expected agent type")
	}

	allowAgentConstructions := map[typegraph.TypeReference]bool{}

	if sc.allowAgentConstructions != nil {
		existing := *sc.allowAgentConstructions
		for key, value := range existing {
			allowAgentConstructions[key] = value
		}
	}

	allowAgentConstructions[agentType] = true

	return scopeContext{
		rootNode:                   sc.rootNode,
		staticDependencyCollector:  sc.staticDependencyCollector,
		dynamicDependencyCollector: sc.dynamicDependencyCollector,

		accessOption:      sc.accessOption,
		overrideTypes:     sc.overrideTypes,
		parentImplemented: sc.parentImplemented,
		parentBreakable:   sc.parentBreakable,
		parentContinuable: sc.parentContinuable,

		rootLabelSet:         sc.rootLabelSet,
		localScopeNamesCache: sc.localScopeNamesCache,

		allowAgentConstructions: &allowAgentConstructions,
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

		accessOption:            sc.accessOption,
		allowAgentConstructions: sc.allowAgentConstructions,
		parentImplemented:       sc.parentImplemented,
		parentBreakable:         sc.parentBreakable,
		parentContinuable:       sc.parentContinuable,

		rootLabelSet:         sc.rootLabelSet,
		localScopeNamesCache: sc.localScopeNamesCache,

		overrideTypes: &overrideTypes,
	}
}
