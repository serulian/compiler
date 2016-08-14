// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"path"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"
)

// IsSamePackage returns true if the given input source paths are found *directly* under the same
// package (no subpackages).
func InSamePackage(first compilercommon.InputSource, second compilercommon.InputSource) bool {
	return first == second || path.Dir(string(first)) == path.Dir(string(second))
}

// getPackageForImport returns the package information for the package imported by the given import
// package node.
func (g *SRG) getPackageForImport(importPackageNode compilergraph.GraphNode) importedPackage {
	importNode := importPackageNode.GetIncomingNode(parser.NodeImportPredicatePackageRef)

	// Note: There may not be a kind, in which case this will return empty string, which is the
	// default kind.
	packageKind, _ := importNode.TryGet(parser.NodeImportPredicateKind)
	packageLocation := importNode.Get(parser.NodeImportPredicateLocation)

	packageInfo, ok := g.packageMap.Get(packageKind, packageLocation)
	if !ok {
		source := importNode.Get(parser.NodeImportPredicateSource)
		subsource, _ := importPackageNode.TryGet(parser.NodeImportPredicateSubsource)
		panic(fmt.Sprintf("Missing package info for import %s %s (reference %v) (node %v)\nPackage Map: %v",
			source, subsource, packageLocation, importNode, g.packageMap))
	}

	return importedPackage{
		srg:          g,
		packageInfo:  packageInfo,
		importSource: compilercommon.InputSource(importPackageNode.Get(parser.NodePredicateSource)),
	}
}

// srgPackage implements the typeContainer information for searching over a package of
// modules.
type importedPackage struct {
	srg          *SRG                       // The parent SRG.
	packageInfo  packageloader.PackageInfo  // The package info for this package.
	importSource compilercommon.InputSource // The input source for the import.
}

// IsSRGPackage returns true if the imported package is an SRG package.
func (p importedPackage) IsSRGPackage() bool {
	return p.packageInfo.Kind() == srgSourceKind
}

// ModulePaths returns the paths of all the modules under this package.
func (p importedPackage) ModulePaths() []compilercommon.InputSource {
	return p.packageInfo.ModulePaths()
}

// ResolveType will attempt to resolve the given type name under all modules in this package.
func (p importedPackage) ResolveType(name string) (TypeResolutionResult, bool) {
	if !p.IsSRGPackage() {
		return resultForExternalPackage(name, p.packageInfo), true
	}

	for _, modulePath := range p.packageInfo.ModulePaths() {
		module, ok := p.srg.FindModuleBySource(modulePath)
		if !ok {
			panic(fmt.Sprintf("Could not find module with path: %s", modulePath))
		}

		result, ok := module.ResolveType(name, p.moduleResolutionOption(modulePath))
		if ok {
			return result, true
		}
	}

	return TypeResolutionResult{}, false
}

// FindTypeOrMemberByName searches all of the modules in this package for a type or member with the given name.
// Will panic for non-SRG imported packages.
func (p importedPackage) FindTypeOrMemberByName(name string) (SRGTypeOrMember, bool) {
	if !p.IsSRGPackage() {
		panic("Cannot call FindTypeOrMemberByName on non-SRG package")
	}

	for _, modulePath := range p.packageInfo.ModulePaths() {
		module, ok := p.srg.FindModuleBySource(modulePath)
		if !ok {
			panic(fmt.Sprintf("Could not find module with path: %s", modulePath))
		}

		namedFound, ok := module.FindTypeOrMemberByName(name, p.moduleResolutionOption(modulePath))
		if ok {
			return namedFound, true
		}
	}

	return SRGTypeOrMember{}, false
}

// moduleResolutionOption returns the resolution option to use when resolving under the given
// module path. If the module is under the same package as the package that created this import,
// then we allow resolution of otherwise "unexported" members and types.
func (p importedPackage) moduleResolutionOption(modulePath compilercommon.InputSource) ModuleResolutionOption {
	if InSamePackage(modulePath, p.importSource) {
		return ModuleResolveAll
	}

	return ModuleResolveExportedOnly
}
