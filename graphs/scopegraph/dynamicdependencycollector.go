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
type dynamicDependencyCollector interface {
	// registerDependency registers the given name as a dynamic dependency.
	registerDynamicDependency(name string)

	// NameSlice returns the dynamic dependencies collected into a slice of strings.
	NameSlice() []string
}

// noopDynamicDependencyCollector defines a dynamic dependency collector which does nothing.
type noopDynamicDependencyCollector struct{}

func (dc *noopDynamicDependencyCollector) registerDynamicDependency(name string) {}
func (dc *noopDynamicDependencyCollector) NameSlice() []string {
	panic("Should never be called")
}

// newDynamicDependencyCollector creates a new dynamic dependency collector. Should be
// called when constructing the scopeContext for scoping a type or module member.
func newDynamicDependencyCollector() dynamicDependencyCollector {
	return &workingDynamicDependencyCollector{
		dependencies: map[string]bool{},
	}
}

// workingDynamicDependencyCollector defines a working dynamic dependency collector.
type workingDynamicDependencyCollector struct {
	dependencies map[string]bool
}

func (dc *workingDynamicDependencyCollector) registerDynamicDependency(name string) {
	dc.dependencies[name] = true
}

func (dc *workingDynamicDependencyCollector) NameSlice() []string {
	dynamicDepSlice := make([]string, len(dc.dependencies))
	index := 0

	for dynamicDependencyName := range dc.dependencies {
		dynamicDepSlice[index] = dynamicDependencyName
		index = index + 1
	}

	return dynamicDepSlice
}
