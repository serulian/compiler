// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"path"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

type ModuleResolutionOption int

const (
	ModuleResolveExportedOnly ModuleResolutionOption = iota
	ModuleResolveAll          ModuleResolutionOption = iota
)

// SRGModule wraps a module defined in the SRG.
type SRGModule struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// GetModules returns all the modules defined in the SRG.
func (g *SRG) GetModules() []SRGModule {
	it := g.findAllNodes(parser.NodeTypeFile).BuildNodeIterator()
	var modules []SRGModule

	for it.Next() {
		modules = append(modules, SRGModule{it.Node(), g})
	}

	return modules
}

// FindModuleBySource returns the module with the given input source, if any.
func (g *SRG) FindModuleBySource(source compilercommon.InputSource) (SRGModule, bool) {
	node, found := g.findAllNodes(parser.NodeTypeFile).
		Has(parser.NodePredicateSource, string(source)).
		TryGetNode()

	return SRGModule{node, g}, found
}

// InputSource returns the input source for this module.
func (m SRGModule) InputSource() compilercommon.InputSource {
	return compilercommon.InputSource(m.GraphNode.Get(parser.NodePredicateSource))
}

// Name returns the name of the module.
func (m SRGModule) Name() string {
	return path.Base(string(m.InputSource()))
}

// Node returns the underlying node.
func (m SRGModule) Node() compilergraph.GraphNode {
	return m.GraphNode
}

// GetTypes returns the types declared directly under the module.
func (m SRGModule) GetTypes() []SRGType {
	it := m.GraphNode.StartQuery().
		Out(parser.NodePredicateChild).
		IsKind(TYPE_KINDS_TAGGED...).
		BuildNodeIterator()

	var types []SRGType

	for it.Next() {
		types = append(types, SRGType{it.Node(), m.srg})
	}

	return types
}

// GetMembers returns the members declared directly under the module.
func (m SRGModule) GetMembers() []SRGMember {
	it := m.GraphNode.StartQuery().
		Out(parser.NodePredicateChild).
		IsKind(parser.NodeTypeFunction, parser.NodeTypeVariable).
		BuildNodeIterator()

	var members []SRGMember

	for it.Next() {
		members = append(members, SRGMember{it.Node(), m.srg})
	}

	return members
}

// findImportByName searches for the import with the given package name and returns it, if any.
func (m SRGModule) findImportByPackageName(name string) (compilergraph.GraphNode, bool) {
	return m.StartQuery().
		Out(parser.NodePredicateChild).
		IsKind(parser.NodeTypeImport).
		Has(parser.NodeImportPredicatePackageName, name).
		TryGetNode()
}

// findImportWithLocalName attempts to find an import in this module with the given local name.
func (m SRGModule) findImportWithLocalName(name string) (compilergraph.GraphNode, bool) {
	return m.GraphNode.StartQuery().
		Out(parser.NodePredicateChild).
		IsKind(parser.NodeTypeImport).
		Has(parser.NodeImportPredicateName, name).
		TryGetNode()
}

// FindTypeOrMemberByName searches for the type definition, declaration or module member with the given
// name under this module and returns it (if found). Note that this method does not handle imports.
func (m SRGModule) FindTypeOrMemberByName(name string, option ModuleResolutionOption) (SRGTypeOrMember, bool) {
	// If only exported members are allowed, ensure that the name being searched matches exported
	// name rules.
	if option == ModuleResolveExportedOnly && !isExportedName(name) {
		return SRGTypeOrMember{}, false
	}

	// Find matching types and members in the module.
	nodeFound, found := m.StartQuery().
		Out(parser.NodePredicateChild).
		Has("named", name).
		IsKind(MEMBER_KINDS_TAGGED...).
		TryGetNode()

	if !found {
		return SRGTypeOrMember{}, false
	}

	// Return the type found.
	return SRGTypeOrMember{nodeFound, m.srg}, true
}

// FindTypeByName searches for the type definition or declaration with the given name under
// this module and returns it (if found).
func (m SRGModule) FindTypeByName(typeName string, option ModuleResolutionOption) (SRGType, bool) {
	typeOrMember, found := m.FindTypeOrMemberByName(typeName, option)
	if !found || !typeOrMember.IsType() {
		return SRGType{}, false
	}

	return SRGType{typeOrMember.GraphNode, m.srg}, true
}
