// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
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

// Subsource returns the subsource for this package import.
func (i SRGPackageImport) Subsource() string {
	return i.GraphNode.Get(parser.NodeImportPredicateSubsource)
}

// ResolvedTypeOrMember returns the SRG type or member referenced by this import, if any.
func (i SRGPackageImport) ResolvedTypeOrMember() (SRGTypeOrMember, bool) {
	// Load the package information.
	packageInfo, err := i.srg.getPackageForImport(i.GraphNode)
	if err != nil || !packageInfo.IsSRGPackage() {
		return SRGTypeOrMember{}, false
	}

	// Search for the subsource.
	return packageInfo.FindTypeOrMemberByName(i.Subsource())
}
