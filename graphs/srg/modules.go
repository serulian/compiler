// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
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
	srg      *SRG                    // The parent SRG.
	fileNode compilergraph.GraphNode // The root file node in the SRG for the module.

	InputSource parser.InputSource // The input source for the module
}

// GetModules returns all the modules defined in the SRG.
func (g *SRG) GetModules() []SRGModule {
	it := g.findAllNodes(parser.NodeTypeFile).BuildNodeIterator(parser.NodePredicateSource)
	var modules []SRGModule

	for it.Next() {
		modules = append(modules, moduleForSRGNode(g, it.Node, it.Values[parser.NodePredicateSource]))
	}

	return modules
}

// FindModuleBySource returns the module with the given input source, if any.
func (g *SRG) FindModuleBySource(source parser.InputSource) (SRGModule, bool) {
	node, found := g.findAllNodes(parser.NodeTypeFile).
		Has(parser.NodePredicateSource, string(source)).
		TryGetNode()

	if !found {
		return SRGModule{}, false
	}

	return moduleForSRGNode(g, node, string(source)), true
}

// FileNode returns the root file node for this module.
func (m SRGModule) FileNode() compilergraph.GraphNode {
	return m.fileNode
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

	typeOrImportNode, found := m.fileNode.StartQuery().
		Out(parser.NodePredicateChild).
		Has(parser.NodeClassPredicateName, typeName).
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
	return typeForSRGNode(m.srg, typeOrImportNode, typeName), true
}

// moduleForSRGNode returns an SRGModule struct representing the node, which is the root
// file node in the SRG for the module.
func moduleForSRGNode(g *SRG, fileNode compilergraph.GraphNode, inputSource string) SRGModule {
	return SRGModule{
		srg:         g,
		fileNode:    fileNode,
		InputSource: parser.InputSource(inputSource),
	}
}
