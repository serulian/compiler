// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilercommon

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/emirpasic/gods/trees/redblacktree"
)

// NewPositionMapper returns a pointer to a struct that can be used for cached, faster mapping
// of byte positions in input source files to their associated 0-indexed line number and
// column positions. Note that this struct is *not* safe for concurrent access.
func NewPositionMapper() *PositionMapper {
	return &PositionMapper{
		sources: map[InputSource]SourcePositionMapper{},
	}
}

// PositionMapper defines a helper struct for cached, faster lookup of byte position ->
// (line, column) for a set of source files.
type PositionMapper struct {
	// sources maps from a particular input source to its associated mapper.
	sources map[InputSource]SourcePositionMapper
}

// Map returns the line number and column position of the byte position in the given source,
// if any.
func (pm *PositionMapper) Map(source InputSource, bytePosition int) (int, int, error) {
	if existing, ok := pm.sources[source]; ok {
		return existing.Map(bytePosition)
	}

	mapper, err := newSourcePositionMapper(source)
	if err != nil {
		return -1, -1, err
	}

	pm.sources[source] = mapper
	return mapper.Map(bytePosition)
}

// SourcePositionMapper defines a helper struct for cached, faster lookup of byte position ->
// (line, column) for a specific source file.
type SourcePositionMapper struct {
	rangeTree *redblacktree.Tree
}

// newSourcePositionMapper creates a new source position mapper for the given source file.
func newSourcePositionMapper(source InputSource) (SourcePositionMapper, error) {
	contents, err := ioutil.ReadFile(string(source))
	if err != nil {
		return SourcePositionMapper{}, err
	}

	return CreateSourcePositionMapper(contents), nil
}

type inclusiveRange struct {
	start int
	end   int
}

type lineAndStart struct {
	lineNumber    int
	startPosition int
}

func inclusiveComparator(a, b interface{}) int {
	i1 := a.(inclusiveRange)
	i2 := b.(inclusiveRange)

	if i1.start >= i2.start && i1.end <= i2.end {
		return 0
	}

	diff := i1.start - i2.start

	if diff < 0 {
		return -1
	}
	if diff > 0 {
		return 1
	}
	return 0
}

// CreateSourcePositionMapper returns a source position mapper for the contents of a source file.
func CreateSourcePositionMapper(contents []byte) SourcePositionMapper {
	lines := strings.Split(string(contents), "\n")
	rangeTree := redblacktree.NewWith(inclusiveComparator)

	var currentStart = 0
	for index, line := range lines {
		lineEnd := currentStart + len(line)
		rangeTree.Put(inclusiveRange{currentStart, lineEnd}, lineAndStart{index, currentStart})
		currentStart = lineEnd + 1
	}

	return SourcePositionMapper{rangeTree}
}

// Map returns the line number and column position of the byte position in the mapper's source,
// if any.
func (spm SourcePositionMapper) Map(bytePosition int) (int, int, error) {
	ls, found := spm.rangeTree.Get(inclusiveRange{bytePosition, bytePosition})
	if !found {
		return -1, -1, fmt.Errorf("Unknown position %v in source file", bytePosition)
	}

	las := ls.(lineAndStart)
	return las.lineNumber, bytePosition - las.startPosition, nil
}
