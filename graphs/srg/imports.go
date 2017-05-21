// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// GetImport returns an SRGImport wrapper around the given import node. Will panic if the node
// is not an import node.
func (g *SRG) GetImport(importNode compilergraph.GraphNode) SRGImport {
	if importNode.Kind() != parser.NodeTypeImport {
		panic("Expected import node")
	}

	return SRGImport{importNode, g}
}

// SRGImport represents an import node.
type SRGImport struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// Source returns the source of the import.
func (i SRGImport) Source() string {
	return i.GraphNode.Get(parser.NodeImportPredicateSource)
}

// Code returns a code-like summarization of the import, for human consumption.
func (i SRGImport) Code() string {
	return "import " + i.Source()
}

// SourceLocation returns the source location for this import's module, if applicable.
func (i SRGImport) SourceLocation() (compilercommon.SourceAndLocation, bool) {
	packageLocation := i.Get(parser.NodeImportPredicateLocation)
	packageKind, _ := i.TryGet(parser.NodeImportPredicateKind)
	packageInfo, ok := i.srg.packageMap.Get(packageKind, packageLocation)
	if !ok {
		return compilercommon.SourceAndLocation{}, false
	}

	modulePaths := packageInfo.ModulePaths()
	if len(modulePaths) == 1 {
		return compilercommon.NewSourceAndLocation(modulePaths[0], 0), true
	}

	return compilercommon.SourceAndLocation{}, false
}

// PackageImports returns the package imports for this import statement, if any.
func (i SRGImport) PackageImports() []SRGPackageImport {
	pit := i.GraphNode.
		StartQuery().
		Out(parser.NodeImportPredicatePackageRef).
		BuildNodeIterator()

	var packageImports = make([]SRGPackageImport, 0)
	for pit.Next() {
		packageImports = append(packageImports, i.srg.GetPackageImport(pit.Node()))
	}
	return packageImports
}

// GetPackageImport returns an SRGPackageImport wrapper around the given import package node.
// Will panic if the node is not an import package node.
func (g *SRG) GetPackageImport(packageNode compilergraph.GraphNode) SRGPackageImport {
	if packageNode.Kind() != parser.NodeTypeImportPackage {
		panic("Expected import package node")
	}

	return SRGPackageImport{packageNode, g}
}

// SRGPackageImport represents the package under an import.
type SRGPackageImport struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// Subsource returns the subsource for this package import, if any.
func (i SRGPackageImport) Subsource() (string, bool) {
	return i.GraphNode.TryGet(parser.NodeImportPredicateSubsource)
}

// Alias returns the local alias for this package import, if any.
func (i SRGPackageImport) Alias() (string, bool) {
	return i.GraphNode.TryGet(parser.NodeImportPredicateName)
}

// SourceLocation returns the source location for this package import's module, if applicable.
func (i SRGPackageImport) SourceLocation() (compilercommon.SourceAndLocation, bool) {
	packageInfo, err := i.srg.getPackageForImport(i.GraphNode)
	if err != nil {
		return compilercommon.SourceAndLocation{}, false
	}

	modulePaths := packageInfo.ModulePaths()
	if len(modulePaths) == 1 {
		return compilercommon.NewSourceAndLocation(modulePaths[0], 0), true
	}

	return compilercommon.SourceAndLocation{}, false
}

// ResolvedTypeOrMember returns the SRG type or member referenced by this import, if any.
func (i SRGPackageImport) ResolvedTypeOrMember() (SRGTypeOrMember, bool) {
	// Load the package information.
	packageInfo, err := i.srg.getPackageForImport(i.GraphNode)
	if err != nil || !packageInfo.IsSRGPackage() {
		return SRGTypeOrMember{}, false
	}

	// Search for the subsource.
	subsource, _ := i.Subsource()
	return packageInfo.FindTypeOrMemberByName(subsource)
}

// Code returns a code-like summarization of the import, for human consumption.
func (i SRGPackageImport) Code() string {
	importNode := i.GraphNode.GetIncomingNode(parser.NodeImportPredicatePackageRef)
	importRef := SRGImport{importNode, i.srg}

	subsource, hasSubsource := i.Subsource()
	if hasSubsource {
		return fmt.Sprintf("from %s import %s", subsource, importRef.Source())
	}

	return "import " + importRef.Source()
}
