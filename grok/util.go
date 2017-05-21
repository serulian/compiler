// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
)

type sourceLocationCapable interface {
	SourceLocation() (compilercommon.SourceAndLocation, bool)
}

type sourceLocationsCapable interface {
	SourceLocations() []compilercommon.SourceAndLocation
}

// getSALs returns the source locations found for the given capable instance, if any.
func getSALs(slc sourceLocationCapable) []compilercommon.SourceAndLocation {
	instance, supportsMultiple := slc.(sourceLocationsCapable)
	if supportsMultiple {
		return instance.SourceLocations()
	}

	sal, hasSal := slc.SourceLocation()
	if hasSal {
		return []compilercommon.SourceAndLocation{sal}
	}

	return []compilercommon.SourceAndLocation{}
}

// getSALsForNode returns a source and location for the given SRG node, if any.
func getSALsForNode(node compilergraph.GraphNode) []compilercommon.SourceAndLocation {
	sal := srg.SourceLocationForNode(node)
	return []compilercommon.SourceAndLocation{sal}
}

// getSALsForTypeRef returns a source and location for the type referenced, if any.
func getSALsForTypeRef(typeref typegraph.TypeReference) []compilercommon.SourceAndLocation {
	if !typeref.IsNormal() {
		return []compilercommon.SourceAndLocation{}
	}

	return getSALs(typeref.ReferredType())
}
