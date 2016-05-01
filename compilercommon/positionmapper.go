// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilercommon

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// NewPositionMapper returns a pointer to a struct that can be used for cached, faster mapping
// of byte positions in input source files to their associated 0-indexed line number and
// column positions. Note that this struct is *not* safe for concurrent access.
func NewPositionMapper() *positionMapper {
	return &positionMapper{
		sources: map[InputSource]sourcePositionMapper{},
	}
}

// positionMapper defines a helper struct for cached, faster lookup of byte position ->
// (line, column) for a set of source files.
type positionMapper struct {
	// sources maps from a particular input source to its associated mapper.
	sources map[InputSource]sourcePositionMapper
}

// Map returns the line number and column position of the byte position in the given source,
// if any.
func (pm *positionMapper) Map(source InputSource, bytePosition int) (int, int, error) {
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

// sourcePositionMapper defines a helper struct for cached, faster lookup of byte position ->
// (line, column) for a specific source file.
type sourcePositionMapper struct {
	lineLengths []int // The lengths of each of the lines in the source file.
}

// newSourcePositionMapper creates a new source position mapper for the given source file.
func newSourcePositionMapper(source InputSource) (sourcePositionMapper, error) {
	contents, err := ioutil.ReadFile(string(source))
	if err != nil {
		return sourcePositionMapper{}, err
	}

	lines := strings.Split(string(contents), "\n")
	lineLengths := make([]int, len(lines))

	for index, line := range lines {
		lineLengths[index] = len(line)
	}

	return sourcePositionMapper{lineLengths}, nil
}

// Map returns the line number and column position of the byte position in the mapper's source,
// if any.
func (spm sourcePositionMapper) Map(bytePosition int) (int, int, error) {
	var currentPosition = 0
	for index, lineLength := range spm.lineLengths {
		nextPosition := currentPosition + lineLength + 1
		if bytePosition < nextPosition {
			return index, bytePosition - currentPosition, nil
		}

		if bytePosition == nextPosition {
			return index + 1, 0, nil
		}

		currentPosition = nextPosition
	}

	return -1, -1, fmt.Errorf("Unknown position %v in source file", bytePosition)
}
