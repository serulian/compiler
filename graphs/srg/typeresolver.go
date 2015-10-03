// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"strings"

	"github.com/serulian/compiler/parser"
)

// typeContainer is an interface that implements the FindTypeByName method.
type typeContainer interface {
	// FindTypeByName searches for the type definition with the given name and returns
	// it if found.
	FindTypeByName(typeName string, option ModuleResolutionOption) (SRGType, bool)
}

// ResolveType attempts to resolve the type at the given path under this module and its
// imports.
func (m SRGModule) ResolveType(path string) (SRGType, bool) {
	pieces := strings.Split(path, ".")

	if len(pieces) < 1 || len(pieces) > 2 {
		panic(fmt.Sprintf("Expected type string with one or two pieces, found: %v", pieces))
	}

	// If there is only a single piece, this is a local-module type or alias.
	if len(pieces) == 1 {
		// Check the aliases map.
		if aliasedType, ok := m.srg.ResolveAliasedType(path); ok {
			return aliasedType, true
		}

		return m.FindTypeByName(path, ModuleResolveAll)
	}

	// Otherwise, we first need to find a type container.
	container, found := m.ResolveTypeContainer(pieces[0])
	if !found {
		return SRGType{}, false
	}

	// Resolve the type under the container.
	return container.FindTypeByName(pieces[1], ModuleResolveExportedOnly)
}

// ResolveTypeContainer attempts to resolve a name under the imports of a module and
// returns the typeContainer found (if any).
//
// Names checked:
//   - Names of imports: `import something`
//   - Aliases of local imports: `import something as SomethingElse`
//   - Aliases of SCM imports: `import "some.com/foo/bar" as SomethingElse`
func (m SRGModule) ResolveTypeContainer(name string) (typeContainer, bool) {
	packageInfo, ok := m.ResolveImportedPackage(name)
	if !ok {
		return nil, false
	}

	// If the package contains a single module file, then return the module itself as the type container.
	if len(packageInfo.ModulePaths()) == 1 {
		module, ok := m.srg.FindModuleBySource(packageInfo.ModulePaths()[0])
		if !ok {
			panic(fmt.Sprintf("Could not find module with path: %s", packageInfo.ModulePaths()[0]))
		}
		return module, true
	}

	// Otherwise, we return an srgPackage helper that searches all the files within.
	return packageInfo, true
}

// ResolveImportedPackage attempts to find a package imported by this module under the given
// packageName.
func (m SRGModule) ResolveImportedPackage(packageName string) (*srgPackage, bool) {
	// Search for the import under the module with the given package name.
	node, found := m.StartQuery().
		Out(parser.NodePredicateChild).
		IsKind(parser.NodeTypeImport).
		Has(parser.NodeImportPredicatePackageName, packageName).
		TryGetNode()

	if !found {
		return nil, false
	}

	// If we've found a node, retrieve its package location and find it
	// in the package map.
	return m.srg.getPackageForImport(node), true
}
