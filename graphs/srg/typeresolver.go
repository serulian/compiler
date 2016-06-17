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

// typeContainer is an interface that implements the FindTypeByName method.
type typeContainer interface {
	// ResolveExportedType searches for the given type path under the container.
	ResolveExportedType(path string) (TypeResolutionResult, bool)
}

type TypeResolutionResult struct {
	IsExternalPackage       bool
	ResolvedType            SRGTypeOrGeneric
	ExternalPackage         packageloader.PackageInfo
	ExternalPackageTypePath string
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

// ResolveExportedType attempts to resolve the type exported.
func (m SRGModule) ResolveExportedType(path string) (TypeResolutionResult, bool) {
	localType, found := m.FindTypeByName(path, ModuleResolveExportedOnly)
	if found {
		return resultForType(localType), true
	}

	return TypeResolutionResult{}, false
}

// ResolveType attempts to resolve the type at the given path under this module and its
// imports.
func (m SRGModule) ResolveType(path string) (TypeResolutionResult, bool) {
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
		packageInfo := m.srg.getPackageForImport(importPackageNode)
		return packageInfo.ResolveExportedType(importPackageNode.Get(parser.NodeImportPredicateSubsource))
	}

	// Otherwise, we first need to find a type container.
	container, found := m.resolveTypeContainer(pieces[0])
	if !found {
		return TypeResolutionResult{}, false
	}

	// Resolve the type under the container.
	return container.ResolveExportedType(pieces[1])
}

// findImportPackageNode searches for the import package node under this module with the given
// matching name found on the given predicate.
func (m SRGModule) findImportPackageNode(name string, predicate string) (compilergraph.GraphNode, bool) {
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
	return m.srg.getPackageForImport(node), true
}

// resolveTypeContainer attempts to resolve a name under the imports of a module and
// returns the typeContainer found (if any).
//
// Names checked:
//   - Names of imports: `import something`
//   - Aliases of local imports: `import something as SomethingElse`
//   - Aliases of SCM imports: `import "some.com/foo/bar" as SomethingElse`
func (m SRGModule) resolveTypeContainer(name string) (typeContainer, bool) {
	packageInfo, ok := m.ResolveImportedPackage(name)
	if !ok {
		return nil, false
	}

	// If the package contains a single SRG module file, then return the module itself as the type container.
	singleModule, hasSingleModule := packageInfo.SingleModule()
	if hasSingleModule {
		return singleModule, true
	}

	// Otherwise, we return an importedPackage helper that searches all the files within.
	return packageInfo, true
}
