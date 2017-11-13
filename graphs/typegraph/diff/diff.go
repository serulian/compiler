// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"strings"

	"github.com/serulian/compiler/graphs/typegraph"
)

// ComputeDiff computes the differences between the original type graph and the
// updated type graph.
func ComputeDiff(original *typegraph.TypeGraph, updated *typegraph.TypeGraph, packageFilters ...string) TypeGraphDiff {
	// Find all the packages in both type graphs and compare.
	originalPackages := computePackages(original, packageFilters)
	updatedPackages := computePackages(updated, packageFilters)

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
		Original: original,
		Updated:  updated,
		Filters:  packageFilters,
	}
}

type packageReference struct {
	source  string
	modules []typegraph.TGModule
}

func computePackages(graph *typegraph.TypeGraph, packageFilters []string) map[string]*packageReference {
	packages := map[string]*packageReference{}
	for _, module := range graph.Modules() {
		packagePath := module.PackagePath()
		if len(packageFilters) > 0 {
			var found = false
			for _, filter := range packageFilters {
				if strings.HasPrefix(packagePath, filter) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if _, found := packages[packagePath]; !found {
			packages[packagePath] = &packageReference{
				source:  packagePath,
				modules: []typegraph.TGModule{},
			}
		}

		packages[packagePath].modules = append(packages[packagePath].modules, module)
	}
	return packages
}
