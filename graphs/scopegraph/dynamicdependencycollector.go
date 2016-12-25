// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"
)

var _ = fmt.Printf

// dynamicDependencyCollector defines a helper for collecting all dynamic dependencies of
// the implementation of a type or module member.
type dynamicDependencyCollector struct {
	dependencies map[string]bool
}

// newDynamicDependencyCollector creates a new dynamic dependency collector. Should be
// called when constructing the scopeContext for scoping a type or module member.
func newDynamicDependencyCollector() *dynamicDependencyCollector {
	return &dynamicDependencyCollector{
		dependencies: map[string]bool{},
	}
}

// registerDependency registers the given name as a dynamic dependency.
func (dc *dynamicDependencyCollector) registerDynamicDependency(name string) {
	dc.dependencies[name] = true
}

// NameSlice returns the dynamic dependencies collected into a slice of strings.
func (dc *dynamicDependencyCollector) NameSlice() []string {
	dynamicDepSlice := make([]string, len(dc.dependencies))
	index := 0

	for dynamicDependencyName, _ := range dc.dependencies {
		dynamicDepSlice[index] = dynamicDependencyName
		index = index + 1
	}

	return dynamicDepSlice
}
