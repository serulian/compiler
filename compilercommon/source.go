// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// compilercommon package defines common types and their functions shared across the compiler.
package compilercommon

import (
	"io/ioutil"
	"path"
	"strings"
)

// InputSource represents the a source file path in the graph.
type InputSource string

// GetLocation returns the SourceLocation for the given byte position in the input source.
func (i InputSource) GetLocation(bytePosition int) SourceLocation {
	contents, err := ioutil.ReadFile(string(i))
	if err != nil {
		return SourceLocation{0, 0, 0}
	}

	return GetSourceLocation(string(contents), bytePosition)
}

// SourceLocation represents a location in an input source file.
type SourceLocation struct {
	lineNumber     int // The 0-indexed line number.
	columnPosition int // The 0-indexed column position.
	bytePosition   int // The 0-indexed byte position.
}

// GetSourceLocation returns the calculated source location for the given byte position in the given
// text.
func GetSourceLocation(text string, bytePosition int) SourceLocation {
	lineNumber := strings.Count(text[:bytePosition], "\n")
	newlineLocation := strings.LastIndex(text[:bytePosition], "\n")
	if newlineLocation < 0 {
		// Since there was no newline, the "synthetic" newline is at position -1
		newlineLocation = -1
	}

	columnPosition := bytePosition - newlineLocation - 1
	return SourceLocation{
		lineNumber:     lineNumber,
		columnPosition: columnPosition,
		bytePosition:   bytePosition,
	}
}

// LineNumber returns the line number (0-indexed) of this location.
func (sl SourceLocation) LineNumber() int {
	return sl.lineNumber
}

// ColumnPosition returns the column position (0-indexed) of this location.
func (sl SourceLocation) ColumnPosition() int {
	return sl.columnPosition
}

// SourceAndLocation contains a source path, as well as a location.
type SourceAndLocation struct {
	source       InputSource // The source file path.
	bytePosition int         // The byte position for this location. Used for lazy-loaded.
}

// Source returns the input source path.
func (sal SourceAndLocation) Source() InputSource {
	return sal.source
}

// SourceName returns the name portion of the source path.
func (sal SourceAndLocation) SourceName() string {
	return path.Base(string(sal.source))
}

// Location returns the SourceLocation for this SourceAndLocation.
func (sal SourceAndLocation) Location() SourceLocation {
	return sal.source.GetLocation(sal.bytePosition)
}

// NewSourceAndLocation returns a new source and location for the given source and byte position.
func NewSourceAndLocation(source InputSource, bytePosition int) SourceAndLocation {
	return SourceAndLocation{
		source:       source,
		bytePosition: bytePosition,
	}
}
