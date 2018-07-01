// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/packageloader"
)

func TestInvalidGrokHandle(t *testing.T) {
	defer compilerutil.DetectGoroutineLeak(t, runtime.NumGoroutine())

	// Load a Grok handle without a valid core library and ensure we get back an error.
	testSourcePath := "tests/basic/basic.seru"
	groker := NewGroker(testSourcePath, []string{}, []packageloader.Library{})
	_, err := groker.GetHandle()
	assert.NotNil(t, err, "Expected error for Grok handle")
}

func TestCanceledGrokHandle(t *testing.T) {
	defer compilerutil.DetectGoroutineLeak(t, runtime.NumGoroutine())

	testSourcePath := "tests/basic/basic.seru"
	groker := NewGroker(testSourcePath, []string{}, []packageloader.Library{packageloader.Library{TESTLIB_PATH, false, "", "testcore"}})

	// Build a few handles in rapid succession, and ensure they leave nothing behind.
	var channels = make([]chan HandleResult, 0)

	for i := 0; i < 5; i++ {
		channels = append(channels, groker.BuildHandle())
	}

	finalChan := groker.BuildHandle()
	result := <-finalChan
	assert.Nil(t, result.Error, "Did not expect error for Grok handle")

	for _, channel := range channels {
		<-channel
	}
}
