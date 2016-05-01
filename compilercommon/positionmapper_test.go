// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilercommon

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPositionMapping(t *testing.T) {
	mapper := NewPositionMapper()
	mappingText, err := ioutil.ReadFile("tests/mapping.txt")

	if !assert.Nil(t, err, "Got error reading mapping file") {
		return
	}

	for bytePosition, _ := range mappingText {
		lineNumber, colPosition, err := mapper.Map(InputSource("tests/mapping.txt"), bytePosition)
		if !assert.Nil(t, err, "Got error mapping file") {
			return
		}

		location := NewSourceAndLocation(InputSource("tests/mapping.txt"), bytePosition).Location()

		// Ensure we get the same result for the individual lookup as the mapped one.
		if !assert.Equal(t, location.LineNumber(), lineNumber, "Line number mismatch for position %v", bytePosition) {
			return
		}

		if !assert.Equal(t, location.ColumnPosition(), colPosition, "Column position mismatch for position %v", bytePosition) {
			return
		}
	}
}
