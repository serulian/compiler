// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPositionMapping(t *testing.T) {
	mappingText, err := ioutil.ReadFile("tests/mapping.txt")
	if !assert.Nil(t, err, "Got error reading mapping file") {
		return
	}

	mapper := createSourcePositionMapper(mappingText)

	for runePosition := range mappingText {
		lineNumber, colPosition, err := mapper.RunePositionToLineAndCol(runePosition)
		if !assert.Nil(t, err, "Got error mapping file") {
			return
		}

		// Check mapping back.
		foundRunePosition, err := mapper.LineAndColToRunePosition(lineNumber, colPosition)
		if !assert.Nil(t, err, "Got error mapping file") {
			return
		}

		if !assert.Equal(t, runePosition, foundRunePosition, "Rune position mismatch for position %v", runePosition) {
			return
		}
	}
}
