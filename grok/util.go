// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/graphs/typegraph"
)

type sourceRangeCapable interface {
	SourceRange() (compilercommon.SourceRange, bool)
}

type multipleSourceRangesCapable interface {
	SourceRanges() []compilercommon.SourceRange
}

// sourceRangesOf returns the source ranges found for the given capable instance, if any.
func sourceRangesOf(src sourceRangeCapable) []compilercommon.SourceRange {
	instance, supportsMultiple := src.(multipleSourceRangesCapable)
	if supportsMultiple {
		return instance.SourceRanges()
	}

	sourceRange, hasSourceRange := src.SourceRange()
	if hasSourceRange {
		return []compilercommon.SourceRange{sourceRange}
	}

	return []compilercommon.SourceRange{}
}

// sourceRangesForTypeRef returns source ranges for the type referenced, if any.
func sourceRangesForTypeRef(typeref typegraph.TypeReference) []compilercommon.SourceRange {
	if !typeref.IsNormal() {
		return []compilercommon.SourceRange{}
	}

	return sourceRangesOf(typeref.ReferredType())
}
