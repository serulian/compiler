// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/parser"
)

// SourceStructureFinder defines a helper type for looking up various structure within the
// source represented by an SRG.
type SourceStructureFinder struct {
	// containingImplementedCache is the cache for the containing implemented node for
	// an SRG node.
	containingImplementedCache *compilerutil.RangeMapTree

	// srg contains the parent SRG.
	srg *SRG
}

// SRGContextScopeName represents a name found in the scope.
type SRGContextScopeName struct {
	compilergraph.GraphNode
	alias string
	srg   *SRG
}

// LocalName returns the locally accessible name of the scope in context.
func (cn SRGContextScopeName) LocalName() (string, bool) {
	if cn.alias != "" {
		return cn.alias, true
	}

	return cn.NamedScope().Name()
}

// NamedScope returns the named scope for this context name.
func (cn SRGContextScopeName) NamedScope() SRGNamedScope {
	return SRGNamedScope{cn.GraphNode, cn.srg}
}

// NewSourceStructureFinder returns a new source structural finder for the current SRG.
func (g *SRG) NewSourceStructureFinder() *SourceStructureFinder {
	return &SourceStructureFinder{
		srg: g,
		containingImplementedCache: compilerutil.NewRangeMapTree(g.calculateContainingImplemented),
	}
}

// TryGetContainingNode returns the containing node of the given node that is one of the given types, if any.
func (f *SourceStructureFinder) TryGetContainingNode(node compilergraph.GraphNode, nodeTypes ...parser.NodeType) (compilergraph.GraphNode, bool) {
	containingFilter := func(q compilergraph.GraphQuery) compilergraph.Query {
		startRune := node.GetValue(parser.NodePredicateStartRune).Int()
		endRune := node.GetValue(parser.NodePredicateEndRune).Int()

		return q.
			HasWhere(parser.NodePredicateStartRune, compilergraph.WhereLTE, startRune).
			HasWhere(parser.NodePredicateEndRune, compilergraph.WhereGTE, endRune)
	}

	return f.srg.findAllNodes(nodeTypes...).
		Has(parser.NodePredicateSource, node.Get(parser.NodePredicateSource)).
		FilterBy(containingFilter).
		TryGetNode()
}

// TryGetContainingModule returns the containing module of the given SRG node, if any.
func (f *SourceStructureFinder) TryGetContainingModule(node compilergraph.GraphNode) (SRGModule, bool) {
	moduleNode, found := f.TryGetContainingNode(node, parser.NodeTypeFile)
	if !found {
		return SRGModule{}, false
	}

	return SRGModule{moduleNode, f.srg}, true
}

// TryGetContainingMemberOrType returns the member or type that contains the given node, if any.
func (f *SourceStructureFinder) TryGetContainingMemberOrType(node compilergraph.GraphNode) (SRGTypeOrMember, bool) {
	memberNode, found := f.TryGetContainingNode(node, TYPE_MEMBER_KINDS...)
	if found {
		return SRGTypeOrMember{memberNode, f.srg}, true
	}

	typeNode, found := f.TryGetContainingNode(node, TYPE_KINDS...)
	if found {
		return SRGTypeOrMember{typeNode, f.srg}, true
	}

	return SRGTypeOrMember{}, false
}

// TryGetContainingType returns the type that contains the given node, if any.
func (f *SourceStructureFinder) TryGetContainingType(node compilergraph.GraphNode) (SRGType, bool) {
	typeNode, found := f.TryGetContainingNode(node, TYPE_KINDS...)
	if !found {
		return SRGType{}, false
	}

	return SRGType{typeNode, f.srg}, true
}

// ContainingImplementedOption defines options for the TryGetContainingImplemented function.
type ContainingImplementedOption string

const (
	// ContainingImplementedInclusive indicates that if the node itself is an implemented, it will be
	// returned.
	ContainingImplementedInclusive = "inclusive"

	// ContainingImplementedExclusive indicates that if the node itself is an implemented, its *containing
	// implemented will be returned (if any).
	ContainingImplementedExclusive = "excluse"
)

// TryGetContainingImplemented returns the member, property or function lambda node that
// contains the given node, if any.
func (f *SourceStructureFinder) TryGetContainingImplemented(node compilergraph.GraphNode) (SRGImplementable, bool) {
	return f.TryGetContainingImplementedWithOption(node, ContainingImplementedInclusive)
}

// TryGetContainingImplementedWithOption returns the containing implemented of the given node with the given option,
// if any.
func (f *SourceStructureFinder) TryGetContainingImplementedWithOption(node compilergraph.GraphNode, option ContainingImplementedOption) (SRGImplementable, bool) {
	startRune := node.GetValue(parser.NodePredicateStartRune).Int()
	endRune := node.GetValue(parser.NodePredicateEndRune).Int()

	source := node.Get(parser.NodePredicateSource)

	if option == ContainingImplementedExclusive {
		startRune = startRune - 1
		endRune = endRune + 1
	}

	runeRange := compilerutil.IntRange{startRune, endRune}

	// Lookup the containing implemented via the cache.
	nodeFound := f.containingImplementedCache.Get(source, runeRange)
	if nodeFound == nil {
		return SRGImplementable{}, false
	}

	return SRGImplementable{nodeFound.(compilergraph.GraphNode), f.srg}, true
}

// ScopeInContext returns all the named scope thats available in the context as defined by the given SRG node.
func (f *SourceStructureFinder) ScopeInContext(node compilergraph.GraphNode) []SRGContextScopeName {
	var namedScopes = make([]SRGContextScopeName, 0)

	// Find the root implement(s) and add variables and/or parameters.
	var currentNode = node
	for {
		parentImplemented, hasParentImplemented := f.TryGetContainingImplemented(currentNode)
		if !hasParentImplemented {
			break
		}

		// Find all variables and values in the range between the parent implemented's start rune and the start rune
		// of the node.
		vit := f.variablesAndValuesUnderContext(
			node.Get(parser.NodePredicateSource),
			parentImplemented.GetValue(parser.NodePredicateStartRune).Int(),
			node.GetValue(parser.NodePredicateStartRune).Int())

		for vit.Next() {
			namedScopes = append(namedScopes, f.importedName(vit.Node()))
		}

		// If the parent implemented is a function, add its parameters.
		for _, parameter := range parentImplemented.Parameters() {
			namedScopes = append(namedScopes, f.importedName(parameter.Node()))
		}

		containingImpl, hasContainingImpl := f.TryGetContainingImplementedWithOption(parentImplemented.Node(), ContainingImplementedExclusive)
		if !hasContainingImpl {
			break
		}

		currentNode = containingImpl.GraphNode
	}

	// Add parent type and/or member generics.
	typeOrMember, hasTypeOrMember := f.TryGetContainingMemberOrType(node)
	if hasTypeOrMember {
		// Add type or member generics.
		for _, generic := range typeOrMember.Generics() {
			namedScopes = append(namedScopes, f.importedName(generic.Node()))
		}

		if _, isType := typeOrMember.AsType(); !isType {
			containingType, hasContainingType := f.TryGetContainingType(node)
			if hasContainingType {
				for _, generic := range containingType.Generics() {
					namedScopes = append(namedScopes, f.importedName(generic.Node()))
				}
			}
		}
	}

	// Find the parent module and add members and imports.
	module, _ := f.TryGetContainingModule(node)

	// Add module members.
	for _, moduleMember := range module.GetMembers() {
		namedScopes = append(namedScopes, f.importedName(moduleMember.GraphNode))
	}

	// Add module types.
	for _, moduleType := range module.GetTypes() {
		namedScopes = append(namedScopes, f.importedName(moduleType.GraphNode))
	}

	// Add any imports found.
	for _, moduleImport := range module.GetImports() {
		packageImports := moduleImport.PackageImports()
		for _, importedItem := range packageImports {
			_, hasSubsource := importedItem.Subsource()
			if hasSubsource {
				typeOrMember, resolved := importedItem.ResolvedTypeOrMember()
				if resolved {
					alias, _ := importedItem.Alias()
					namedScopes = append(namedScopes, f.importedNameWithAlias(typeOrMember.GraphNode, alias))
				}
			} else {
				namedScopes = append(namedScopes, f.importedName(importedItem.GraphNode))
			}
		}
	}

	return namedScopes
}

func (f *SourceStructureFinder) importedName(node compilergraph.GraphNode) SRGContextScopeName {
	return SRGContextScopeName{node, "", f.srg}
}

func (f *SourceStructureFinder) importedNameWithAlias(node compilergraph.GraphNode, alias string) SRGContextScopeName {
	return SRGContextScopeName{node, alias, f.srg}
}

func (f *SourceStructureFinder) variablesAndValuesUnderContext(source string, startRune int, endRune int) compilergraph.NodeIterator {
	containingFilter := func(q compilergraph.GraphQuery) compilergraph.Query {
		return q.
			HasWhere(parser.NodePredicateStartRune, compilergraph.WhereGTE, startRune).
			HasWhere(parser.NodePredicateEndRune, compilergraph.WhereLTE, endRune)
	}

	return f.srg.findAllNodes(parser.NodeTypeVariableStatement, parser.NodeTypeNamedValue).
		Has(parser.NodePredicateSource, source).
		FilterBy(containingFilter).
		BuildNodeIterator()
}
