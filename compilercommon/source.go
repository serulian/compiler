// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package compilercommon defines common types and their functions shared across the compiler.
package compilercommon

// Position represents a position in an arbitrary source file.
type Position struct {
	// LineNumber is the 0-indexed line number.
	LineNumber int

	// ColumnPosition is the 0-indexed column position on the line.
	ColumnPosition int
}

// InputSource represents the path of a source file.
type InputSource string

// RangeForRunePosition returns a source range over this source file.
func (is InputSource) RangeForRunePosition(runePosition int, mapper PositionMapper) SourceRange {
	return is.RangeForRunePositions(runePosition, runePosition, mapper)
}

// PositionForRunePosition returns a source position over this source file.
func (is InputSource) PositionForRunePosition(runePosition int, mapper PositionMapper) SourcePosition {
	return runeIndexedPosition{is, mapper, runePosition}
}

// PositionFromLineAndColumn returns a source position at the given line and column in this source file.
func (is InputSource) PositionFromLineAndColumn(lineNumber int, columnPosition int, mapper PositionMapper) SourcePosition {
	return lcIndexedPosition{is, mapper, Position{lineNumber, columnPosition}}
}

// RangeForRunePositions returns a source range over this source file.
func (is InputSource) RangeForRunePositions(startRune int, endRune int, mapper PositionMapper) SourceRange {
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
	RunePositionToLineAndCol(runePosition int, path InputSource) (int, int, error)

	// LineAndColToRunePosition converts the given 0-indexed line number and column position under the
	// given source file into a 0-indexed rune position.
	LineAndColToRunePosition(lineNumber int, colPosition int, path InputSource) (int, error)

	// TextForLine returns the text for the specified line number.
	TextForLine(lineNumber int, path InputSource) (string, error)
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
	RunePosition() (int, error)

	// LineAndColumn returns the 0-indexed line number and column position in the source file.
	LineAndColumn() (int, int, error)

	// LineText returns the text of the line for this position.
	LineText() (string, error)
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
	runePosition int
}

func (ris runeIndexedPosition) Source() InputSource {
	return ris.source
}

func (ris runeIndexedPosition) RunePosition() (int, error) {
	return ris.runePosition, nil
}

func (ris runeIndexedPosition) LineAndColumn() (int, int, error) {
	return ris.mapper.RunePositionToLineAndCol(ris.runePosition, ris.source)
}

func (ris runeIndexedPosition) LineText() (string, error) {
	lineNumber, _, err := ris.LineAndColumn()
	if err != nil {
		return "", err
	}

	return ris.mapper.TextForLine(lineNumber, ris.source)
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

func (lcip lcIndexedPosition) RunePosition() (int, error) {
	return lcip.mapper.LineAndColToRunePosition(lcip.lcPosition.LineNumber, lcip.lcPosition.ColumnPosition, lcip.source)
}

func (lcip lcIndexedPosition) LineAndColumn() (int, int, error) {
	return lcip.lcPosition.LineNumber, lcip.lcPosition.ColumnPosition, nil
}

func (lcip lcIndexedPosition) LineText() (string, error) {
	return lcip.mapper.TextForLine(lcip.lcPosition.LineNumber, lcip.source)
}
