// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"strings"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"
)

// TypeResolutionResult defines the result of resolving a type or package.
type TypeResolutionResult struct {
	// IsExternalPackage returns whether the resolved type is, in fact, an external package.
	IsExternalPackage       bool
	ResolvedType            SRGTypeOrGeneric
	ExternalPackage         packageloader.PackageInfo
	ExternalPackageTypePath string
}

// ResolveType attempts to resolve the type with the given name, under this module.
func (m SRGModule) ResolveType(path string, option ModuleResolutionOption) (TypeResolutionResult, bool) {
	localType, found := m.FindTypeByName(path, option)
	if found {
		return resultForType(localType), true
	}

	return TypeResolutionResult{}, false
}

type trrCache struct {
	result TypeResolutionResult
	ok     bool
}

// ResolveTypePath attempts to resolve the type at the given path under this module and its
// imports.
func (m SRGModule) ResolveTypePath(path string) (TypeResolutionResult, bool) {
	cacheKey := string(m.InputSource()) + ":" + path
	if foundType, ok := m.srg.moduleTypeCache.Get(cacheKey); ok {
		cached := foundType.(trrCache)
		return cached.result, cached.ok
	}

	resolved, ok := m.resolveTypePathNoCaching(path)
	m.srg.moduleTypeCache.Set(cacheKey, trrCache{resolved, ok})
	return resolved, ok
}

func (m SRGModule) resolveTypePathNoCaching(path string) (TypeResolutionResult, bool) {
	pieces := strings.Split(path, ".")

	if len(pieces) < 1 || len(pieces) > 2 {
		panic(fmt.Sprintf("Expected type string with one or two pieces, found: %v", pieces))
	}

	// If there is only a single piece, this is a local-module type, alias or reference to
	// an import.
	if len(pieces) == 1 {
		// Check the aliases map.
		if aliasedType, ok := m.srg.ResolveAliasedType(path); ok {
			return resultForType(aliasedType), true
		}

		// Check for local types.
		localType, found := m.FindTypeByName(path, ModuleResolveAll)
		if found {
			return resultForType(localType), true
		}

		// Check for an imported item by name.
		importPackageNode, found := m.findImportWithLocalName(path)
		if !found {
			return TypeResolutionResult{}, false
		}

		// Resolve the name as the subsource under the import's package.
		packageInfo, err := m.srg.getPackageForImport(importPackageNode)
		if err != nil {
			return TypeResolutionResult{}, false
		}

		return packageInfo.ResolveType(importPackageNode.Get(parser.NodeImportPredicateSubsource))
	}

	// Otherwise, we first need to find an import.
	packageInfo, found := m.ResolveImportedPackage(pieces[0])
	if !found {
		return TypeResolutionResult{}, false
	}

	// Resolve the type under the package info.
	return packageInfo.ResolveType(pieces[1])
}

// ResolveImportedPackage attempts to find a package imported by this module under the given
// packageName.
func (m SRGModule) ResolveImportedPackage(packageName string) (importedPackage, bool) {
	// Search for the imported package under the module with the given package name.
	node, found := m.findImportByPackageName(packageName)
	if !found {
		return importedPackage{}, false
	}

	// If we've found a node, retrieve its package location and find it
	// in the package map.
	info, err := m.srg.getPackageForImport(node)
	if err != nil {
		return importedPackage{}, false
	}

	return info, true
}

func resultForTypeOrGeneric(srgType SRGTypeOrGeneric) TypeResolutionResult {
	return TypeResolutionResult{
		IsExternalPackage: false,
		ResolvedType:      srgType,
	}
}

func resultForType(srgType SRGType) TypeResolutionResult {
	return TypeResolutionResult{
		IsExternalPackage: false,
		ResolvedType:      SRGTypeOrGeneric{srgType.GraphNode, srgType.srg},
	}
}

func resultForExternalPackage(path string, packageInfo packageloader.PackageInfo) TypeResolutionResult {
	return TypeResolutionResult{
		IsExternalPackage:       true,
		ExternalPackage:         packageInfo,
		ExternalPackageTypePath: path,
	}
}

// findImportPackageNode searches for the import package node under this module with the given
// matching name found on the given predicate.
func (m SRGModule) findImportPackageNode(name string, predicate compilergraph.Predicate) (compilergraph.GraphNode, bool) {
	return m.srg.layer.StartQuery(name).
		In(predicate).
		IsKind(parser.NodeTypeImportPackage).
		Has(parser.NodePredicateSource, string(m.InputSource())).
		TryGetNode()
}

// findImportWithLocalName attempts to find an import package in this module with the given local name.
func (m SRGModule) findImportWithLocalName(name string) (compilergraph.GraphNode, bool) {
	return m.findImportPackageNode(name, parser.NodeImportPredicateName)
}

// findImportByPackageName searches for the import package with the given package name and returns it, if any.
func (m SRGModule) findImportByPackageName(packageName string) (compilergraph.GraphNode, bool) {
	return m.findImportPackageNode(packageName, parser.NodeImportPredicatePackageName)
}
