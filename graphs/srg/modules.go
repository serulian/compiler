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

// findImportedByName searches for the imported module or package with the given name and returns it, if any.
func (m SRGModule) findImportedByName(name string) (compilergraph.GraphNode, bool) {
	importNode, importFound := m.StartQuery().
		Out(parser.NodePredicateChild).
		Has(parser.NodeImportPredicatePackageName, name).
		TryGetNode()

	if importFound {
		packageInfo := m.srg.getPackageForImport(importNode)
		singleModule, hasSingleModule := packageInfo.SingleModule()
		if hasSingleModule {
			return singleModule.GraphNode, true
		}

		// TODO(jschorr): Handle packages with more than a single module.
	}

	return compilergraph.GraphNode{}, false
}

// FindTypeByName searches for the type definition or declaration with the given name under
// this module and returns it (if found).
func (m SRGModule) FindTypeByName(typeName string, option ModuleResolutionOption) (SRGType, bool) {
	// If only exported members are allowed, ensure that the name being searched matches exported
	// name rules.
	if option == ModuleResolveExportedOnly && !isExportedName(typeName) {
		return SRGType{}, false
	}

	// By default, we find matching types in the module. If the resolution option allows for
	// imports as well, add them to the allowed list.
	allowedKinds := []compilergraph.TaggedValue{parser.NodeTypeClass, parser.NodeTypeInterface}
	if option == ModuleResolveAll {
		allowedKinds = append(allowedKinds, parser.NodeTypeImport)
	}

	typeOrImportNode, found := m.StartQuery().
		Out(parser.NodePredicateChild).
		Has(parser.NodeTypeDefinitionName, typeName).
		IsKind(allowedKinds...).
		TryGetNode()

	if !found {
		return SRGType{}, false
	}

	// If the node is an import node, find the type under the imported module/package.
	if typeOrImportNode.Kind == parser.NodeTypeImport {
		packageInfo := m.srg.getPackageForImport(typeOrImportNode)

		// If there is a subsource on the import, then we resolve that type name under the package.
		if subsource, ok := typeOrImportNode.TryGet(parser.NodeImportPredicateSubsource); ok {
			return packageInfo.FindTypeByName(subsource, ModuleResolveExportedOnly)
		}

		// Otherwise, we resolve the original requested type name.
		return packageInfo.FindTypeByName(typeName, ModuleResolveExportedOnly)
	}

	// Otherwise, return the type found.
	return SRGType{typeOrImportNode, m.srg}, true
}
