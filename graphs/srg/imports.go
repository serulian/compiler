// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
	"github.com/serulian/compiler/sourceshape"
)

// GetImport returns an SRGImport wrapper around the given import node. Will panic if the node
// is not an import node.
func (g *SRG) GetImport(importNode compilergraph.GraphNode) SRGImport {
	if importNode.Kind() != sourceshape.NodeTypeImport {
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
func (i SRGImport) Source() (string, bool) {
	return i.GraphNode.TryGet(sourceshape.NodeImportPredicateSource)
}

// ParsedSource returns the source for this import, parsed.
func (i SRGImport) ParsedSource() (string, parser.ParsedImportType, error) {
	source, hasSource := i.Source()
	if !hasSource {
		return "", parser.ParsedImportTypeLocal, fmt.Errorf("Missing source")
	}

	return parser.ParseImportValue(source)
}

// Code returns a code-like summarization of the import, for human consumption.
func (i SRGImport) Code() (compilercommon.CodeSummary, bool) {
	source, hasSource := i.Source()
	if !hasSource {
		return compilercommon.CodeSummary{}, false
	}

	return compilercommon.CodeSummary{"", "import " + source, true}, true
}

// SourceRange returns the source range for this import.
func (i SRGImport) SourceRange() (compilercommon.SourceRange, bool) {
	return i.srg.SourceRangeOf(i.GraphNode)
}

// PackageImports returns the package imports for this import statement, if any.
func (i SRGImport) PackageImports() []SRGPackageImport {
	pit := i.GraphNode.
		StartQuery().
		Out(sourceshape.NodeImportPredicatePackageRef).
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
	if packageNode.Kind() != sourceshape.NodeTypeImportPackage {
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
	return i.GraphNode.TryGet(sourceshape.NodeImportPredicateSubsource)
}

// Alias returns the local alias for this package import, if any.
func (i SRGPackageImport) Alias() (string, bool) {
	return i.GraphNode.TryGet(sourceshape.NodeImportPredicateName)
}

// SourceRange returns the source range for this import.
func (i SRGPackageImport) SourceRange() (compilercommon.SourceRange, bool) {
	return i.srg.SourceRangeOf(i.GraphNode)
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
func (i SRGPackageImport) Code() (compilercommon.CodeSummary, bool) {
	importNode := i.GraphNode.GetIncomingNode(sourceshape.NodeImportPredicatePackageRef)
	importRef := SRGImport{importNode, i.srg}

	source, hasSource := importRef.Source()
	if !hasSource {
		return compilercommon.CodeSummary{}, false
	}

	subsource, hasSubsource := i.Subsource()
	if hasSubsource {
		return compilercommon.CodeSummary{"", fmt.Sprintf("from %s import %s", subsource, source), true}, true
	}

	return compilercommon.CodeSummary{"", "import " + source, true}, true
}
