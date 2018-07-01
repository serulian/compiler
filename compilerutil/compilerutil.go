// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// compilerutil package defines utility methods.
package compilerutil

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"testing"

	uuid "github.com/nu7hatch/gouuid"

	"github.com/stretchr/testify/require"
)

type checkFn func() bool

// DCHECK executes the checker function only if the DEBUG environment variable is set.
// If the function returns false, the compiler will panic with the formatted message.
func DCHECK(checker checkFn, failMessage string, args ...interface{}) {
	if os.Getenv("DEBUG") == "" {
		return
	}

	if !checker() {
		panic(fmt.Sprintf(failMessage, args...))
	}
}

// ParseFloat parses the given string value into a float64, panicing on failure.
func ParseFloat(strValue string) float64 {
	value, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		panic(fmt.Sprintf("Expected numeric value, got: %v", value))
	}
	return value
}

// IsId returns whether the given string is a possible ID as returned by NewUniqueId.
func IsId(possibleId string) bool {
	_, err := uuid.ParseHex(possibleId)
	return err == nil
}

// NewUniqueId returns a new unique ID.
func NewUniqueId() string {
	u4, err := uuid.NewV4()
	if err != nil {
		panic(err)
		return ""
	}

	return u4.String()
}

// DetectGoroutineLeak detects if there is a leak of goroutines over the life of a test.
func DetectGoroutineLeak(t *testing.T, grCount int) {
	runtime.GC()
	time.Sleep(1 * time.Millisecond)
	buf := make([]byte, 1<<20)
	runtime.Stack(buf, true)
	require.Equal(t, grCount, runtime.NumGoroutine(), "wrong number of goroutines:\n%s", string(buf))
}
