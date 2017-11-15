// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"path/filepath"

	"github.com/serulian/compiler/graphs/typegraph"
)

type TypeGraphInformation struct {
	Graph           *typegraph.TypeGraph
	PackageRootPath string
}

type DiffFilter func(module typegraph.TGModule) bool

// ComputeDiff computes the differences between the original type graph and the
// updated type graph.
func ComputeDiff(original TypeGraphInformation, updated TypeGraphInformation, filters ...DiffFilter) TypeGraphDiff {
	// Find all the packages in both type graphs and compare.
	originalPackages := computePackages(original, filters)
	updatedPackages := computePackages(updated, filters)

	diffs := map[string]PackageDiff{}

	for path, originalPackage := range originalPackages {
		updatedPackage, hasUpdatedPackage := updatedPackages[path]
		if !hasUpdatedPackage {
			diffs[path] = PackageDiff{
				Kind:            Removed,
				Path:            path,
				ChangeReason:    PackageDiffReasonNotApplicable,
				OriginalModules: originalPackage.modules,
				UpdatedModules:  []typegraph.TGModule{},
				Types:           []TypeDiff{},
				Members:         []MemberDiff{},
			}
			continue
		}

		diffs[path] = diffPackage(path, originalPackage.modules, updatedPackage.modules)
	}

	for path, updatedPackage := range updatedPackages {
		_, hasOriginalPackage := originalPackages[path]
		if !hasOriginalPackage {
			diffs[path] = PackageDiff{
				Kind:            Added,
				Path:            path,
				ChangeReason:    PackageDiffReasonNotApplicable,
				OriginalModules: []typegraph.TGModule{},
				UpdatedModules:  updatedPackage.modules,
				Types:           []TypeDiff{},
				Members:         []MemberDiff{},
			}
			continue
		}
	}

	return TypeGraphDiff{
		Packages: diffs,
		Original: original.Graph,
		Updated:  updated.Graph,
	}
}

type packageReference struct {
	source  string
	modules []typegraph.TGModule
}

func computePackages(graphInfo TypeGraphInformation, filters []DiffFilter) map[string]*packageReference {
	packages := map[string]*packageReference{}
	for _, module := range graphInfo.Graph.Modules() {
		if len(filters) > 0 {
			var found = false
			for _, filter := range filters {
				if filter(module) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		packagePath := module.PackagePath()
		relativePath, err := filepath.Rel(graphInfo.PackageRootPath, packagePath)
		if err != nil {
			continue
		}

		if _, found := packages[relativePath]; !found {
			packages[relativePath] = &packageReference{
				source:  relativePath,
				modules: []typegraph.TGModule{},
			}
		}

		packages[relativePath].modules = append(packages[relativePath].modules, module)
	}
	return packages
}
