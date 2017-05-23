// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"fmt"
	"strings"

	"github.com/emirpasic/gods/trees/redblacktree"
)

// sourcePositionMapper defines a helper struct for cached, faster lookup of rune position <->
// (line, column) for a specific source file.
type sourcePositionMapper struct {
	// rangeTree holds a tree that maps from rune position to a line and start position.
	rangeTree *redblacktree.Tree

	// lineMap holds a map from line number to rune positions for that line.
	lineMap map[uint64]inclusiveRange
}

type inclusiveRange struct {
	start uint64
	end   uint64
}

type lineAndStart struct {
	lineNumber    uint64
	startPosition uint64
}

func inclusiveComparator(a, b interface{}) int {
	i1 := a.(inclusiveRange)
	i2 := b.(inclusiveRange)

	if i1.start >= i2.start && i1.end <= i2.end {
		return 0
	}

	diff := int64(i1.start) - int64(i2.start)

	if diff < 0 {
		return -1
	}
	if diff > 0 {
		return 1
	}
	return 0
}

// createSourcePositionMapper returns a source position mapper for the contents of a source file.
func createSourcePositionMapper(contents []byte) *sourcePositionMapper {
	lines := strings.Split(string(contents), "\n")
	rangeTree := redblacktree.NewWith(inclusiveComparator)
	lineMap := map[uint64]inclusiveRange{}

	var currentStart = uint64(0)
	for index, line := range lines {
		lineEnd := currentStart + uint64(len(line))
		rangeTree.Put(inclusiveRange{currentStart, lineEnd}, lineAndStart{uint64(index), currentStart})
		lineMap[uint64(index)] = inclusiveRange{currentStart, lineEnd}
		currentStart = lineEnd + 1
	}

	return &sourcePositionMapper{rangeTree, lineMap}
}

// RunePositionToLineAndCol returns the line number and column position of the rune position in source.
func (spm *sourcePositionMapper) RunePositionToLineAndCol(runePosition uint64) (uint64, uint64, error) {
	ls, found := spm.rangeTree.Get(inclusiveRange{runePosition, runePosition})
	if !found {
		return 0, 0, fmt.Errorf("Unknown rune position %v in source file", runePosition)
	}

	las := ls.(lineAndStart)
	return las.lineNumber, runePosition - las.startPosition, nil
}

// LineAndColToRunePosition returns the rune position of the line number and column position in source.
func (spm *sourcePositionMapper) LineAndColToRunePosition(lineNumber uint64, colPosition uint64) (uint64, error) {
	lineRuneInfo, hasLine := spm.lineMap[lineNumber]
	if !hasLine {
		return 0, fmt.Errorf("Unknown line %v in source file", lineNumber)
	}

	if colPosition > lineRuneInfo.end-lineRuneInfo.start {
		return 0, fmt.Errorf("Column position %v not found on line %v in source file", colPosition, lineNumber)
	}

	return lineRuneInfo.start + colPosition, nil
}
