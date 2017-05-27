// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilercommon

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// LocalFilePositionMapper is a struct which implements the PositionMapper interface over the local file
// system. Note that since accesses in this struct are *not* cached, this instance should *only* be used
// for testing.
type LocalFilePositionMapper struct{}

func (pm LocalFilePositionMapper) RunePositionToLineAndCol(runePosition int, path InputSource, sourceOption SourceMappingOption) (int, int, error) {
	contents, err := ioutil.ReadFile(string(path))
	if err != nil {
		return -1, -1, err
	}

	sm := CreateSourcePositionMapper(contents)
	return sm.RunePositionToLineAndCol(runePosition)
}

func (pm LocalFilePositionMapper) LineAndColToRunePosition(lineNumber int, colPosition int, path InputSource, sourceOption SourceMappingOption) (int, error) {
	contents, err := ioutil.ReadFile(string(path))
	if err != nil {
		return -1, err
	}

	sm := CreateSourcePositionMapper(contents)
	return sm.LineAndColToRunePosition(lineNumber, colPosition)
}

func (pm LocalFilePositionMapper) TextForLine(lineNumber int, path InputSource, sourceOption SourceMappingOption) (string, error) {
	contents, err := ioutil.ReadFile(string(path))
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(contents), "\n")
	if lineNumber >= len(lines) {
		return "", fmt.Errorf("Invalid line number")
	}

	return lines[lineNumber], nil
}
