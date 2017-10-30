// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/serulian/compiler/packageloader"
)

func TestInvalidGrokHandle(t *testing.T) {
	// Load a Grok handle without a valid core library and ensure we get back an error.
	testSourcePath := "tests/basic/basic.seru"
	groker := NewGroker(testSourcePath, []string{}, []packageloader.Library{})
	_, err := groker.GetHandle()
	assert.NotNil(t, err, "Expected error for Grok handle")
}
