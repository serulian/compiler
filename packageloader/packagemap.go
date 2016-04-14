// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"fmt"

	"github.com/streamrail/concurrent-map"
)

// packageKey defines a unique key string for a source kind and package path.
func packageKey(sourceKind string, packagePath string) string {
	return sourceKind + "::" + packagePath
}

// mutablePackageMap defines a helper type for constructing the package map in a
// concurrent safe way. This is defined to encapsulate a ConcurrentMap in a slightly
// nicer interface.
type mutablePackageMap struct {
	internalMap cmap.ConcurrentMap // The internal concurrent map from packageKey -> PackageInfo.
}

// newMutablePackageMap returns a new mutablePackageMap.
func newMutablePackageMap() *mutablePackageMap {
	return &mutablePackageMap{
		internalMap: cmap.New(),
	}
}

// Add adds a mapping of a package with a specific kind and path to the given info. Will panic
// if a package with the same kind and path has already been added.
func (mpm *mutablePackageMap) Add(sourceKind string, packagePath string, info PackageInfo) {
	if !mpm.internalMap.SetIfAbsent(packageKey(sourceKind, packagePath), info) {
		panic(fmt.Sprintf("Package %s:%s was already added", sourceKind, packagePath))
	}
}

// Build builds the mutable package map into an immutable view.
func (mpm *mutablePackageMap) Build() LoadedPackageMap {
	packageMap := map[string]PackageInfo{}
	for entry := range mpm.internalMap.Iter() {
		packageMap[entry.Key] = entry.Val.(PackageInfo)
	}

	return LoadedPackageMap{packageMap}
}

// LoadedPackageMap defines an immutable view of the loaded packages.
type LoadedPackageMap struct {
	packageMap map[string]PackageInfo
}

// Get returns the package information for the package of the given kind and path, if any.
func (lpm LoadedPackageMap) Get(sourceKind string, packagePath string) (PackageInfo, bool) {
	info, ok := lpm.packageMap[packageKey(sourceKind, packagePath)]
	return info, ok
}
