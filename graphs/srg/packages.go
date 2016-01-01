// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"
)

// getPackageForImport returns the package information for the package imported by the given import
// node.
func (g *SRG) getPackageForImport(importNode compilergraph.GraphNode) *srgPackage {
	packageReference := importNode.Get(parser.NodeImportPredicateLocation)
	packageInfo, ok := g.packageMap[packageReference]
	if !ok {
		source := importNode.Get(parser.NodeImportPredicateSource)
		subsource, _ := importNode.TryGet(parser.NodeImportPredicateSubsource)

		panic(fmt.Sprintf("Missing package info for import %s %s (reference %v) (node %v)\nPackage Map: %v",
			source, subsource, packageReference, importNode, g.packageMap))
	}

	return &srgPackage{
		srg:         g,
		packageInfo: packageInfo,
	}
}

// srgPackage implements the typeContainer information for searching over a package of
// modules.
type srgPackage struct {
	srg         *SRG                      // The parent SRG.
	packageInfo packageloader.PackageInfo // The package info for this package.
}

// ModulePaths returns the paths of all the modules under this package.
func (p *srgPackage) ModulePaths() []compilercommon.InputSource {
	return p.packageInfo.ModulePaths()
}

// SingleModule returns the single module in this package, if any.
func (p *srgPackage) SingleModule() (SRGModule, bool) {
	if len(p.packageInfo.ModulePaths()) == 1 {
		modulePath := p.packageInfo.ModulePaths()[0]
		module, ok := p.srg.FindModuleBySource(modulePath)
		if !ok {
			panic(fmt.Sprintf("Could not find module with path: %s", modulePath))
		}

		return module, true
	}

	return SRGModule{}, false
}

// FindTypeByName searches all of the modules in this package for a type with the given name.
func (p *srgPackage) FindTypeByName(typeName string, option ModuleResolutionOption) (SRGType, bool) {
	for _, modulePath := range p.packageInfo.ModulePaths() {
		module, ok := p.srg.FindModuleBySource(modulePath)
		if !ok {
			panic(fmt.Sprintf("Could not find module with path: %s", modulePath))
		}

		namedFound, ok := module.FindTypeByName(typeName, option)
		if ok {
			return namedFound, true
		}
	}

	return SRGType{}, false
}

// FindTypeOrMemberByName searches all of the modules in this package for a type or member with the given name.
func (p *srgPackage) FindTypeOrMemberByName(name string, option ModuleResolutionOption) (SRGTypeOrMember, bool) {
	for _, modulePath := range p.packageInfo.ModulePaths() {
		module, ok := p.srg.FindModuleBySource(modulePath)
		if !ok {
			panic(fmt.Sprintf("Could not find module with path: %s", modulePath))
		}

		namedFound, ok := module.FindTypeOrMemberByName(name, option)
		if ok {
			return namedFound, true
		}
	}

	return SRGTypeOrMember{}, false
}
