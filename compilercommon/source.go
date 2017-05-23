// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package compilercommon defines common types and their functions shared across the compiler.
package compilercommon

// Position represents a position in an arbitrary source file.
type Position struct {
	// LineNumber is the 0-indexed line number.
	LineNumber uint64

	// ColumnPosition is the 0-indexed column position on the line.
	ColumnPosition uint64
}

// InputSource represents the path of a source file.
type InputSource string

// RangeForRunePositions returns a source range over this source file.
func (is InputSource) RangeForRunePositions(startRune uint64, endRune uint64, mapper PositionMapper) SourceRange {
	return sourceRange{is, runeIndexedPosition{is, mapper, startRune}, runeIndexedPosition{is, mapper, endRune}}
}

// RangeForLineAndColPositions returns a source range over this source file.
func (is InputSource) RangeForLineAndColPositions(start Position, end Position, mapper PositionMapper) SourceRange {
	return sourceRange{is, lcIndexedPosition{is, mapper, start}, lcIndexedPosition{is, mapper, end}}
}

// PositionMapper defines an interface for converting rune position <-> line+col position
// under source files.
type PositionMapper interface {
	// RunePositionToLineAndCol converts the given 0-indexed rune position under the given source file
	// into a 0-indexed line number and column position.
	RunePositionToLineAndCol(runePosition uint64, path InputSource) (uint64, uint64, error)

	// LineAndColToRunePosition converts the given 0-indexed line number and column position under the
	// given source file into a 0-indexed rune position.
	LineAndColToRunePosition(lineNumber uint64, colPosition uint64, path InputSource) (uint64, error)
}

// SourceRange represents a range inside a source file.
type SourceRange interface {
	// Source is the input source for this range.
	Source() InputSource

	// Start is the starting position of the source range.
	Start() SourcePosition

	// End is the ending position (inclusive) of the source range. If the same as the Start,
	// this range represents a single position.
	End() SourcePosition
}

// SourcePosition represents a single position in a source file.
type SourcePosition interface {
	// Source is the input source for this position.
	Source() InputSource

	// RunePosition returns the 0-indexed rune position in the source file.
	RunePosition() (uint64, error)

	// LineAndColumn returns the 0-indexed line number and column position in the source file.
	LineAndColumn() (uint64, uint64, error)
}

// sourceRange implements the SourceRange interface.
type sourceRange struct {
	source InputSource
	start  SourcePosition
	end    SourcePosition
}

func (sr sourceRange) Source() InputSource {
	return sr.source
}

func (sr sourceRange) Start() SourcePosition {
	return sr.start
}

func (sr sourceRange) End() SourcePosition {
	return sr.end
}

// runeIndexedPosition implements the SourcePosition interface over a rune position.
type runeIndexedPosition struct {
	source       InputSource
	mapper       PositionMapper
	runePosition uint64
}

func (ris runeIndexedPosition) Source() InputSource {
	return ris.source
}

func (ris runeIndexedPosition) RunePosition() (uint64, error) {
	return ris.runePosition, nil
}

func (ris runeIndexedPosition) LineAndColumn() (uint64, uint64, error) {
	return ris.mapper.RunePositionToLineAndCol(ris.runePosition, ris.source)
}

// lcIndexedPosition implements the SourcePosition interface over a line and colu,n position.
type lcIndexedPosition struct {
	source     InputSource
	mapper     PositionMapper
	lcPosition Position
}

func (lcip lcIndexedPosition) Source() InputSource {
	return lcip.source
}

func (lcip lcIndexedPosition) RunePosition() (uint64, error) {
	return lcip.mapper.LineAndColToRunePosition(lcip.lcPosition.LineNumber, lcip.lcPosition.ColumnPosition, lcip.source)
}

func (lcip lcIndexedPosition) LineAndColumn() (uint64, uint64, error) {
	return lcip.lcPosition.LineNumber, lcip.lcPosition.ColumnPosition, nil
}
