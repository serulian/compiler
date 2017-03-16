// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilercommon

import (
	"fmt"
)

// SourceError represents an error produced by a source file at a specific location.
type SourceError struct {
	message string            // The error message.
	sal     SourceAndLocation // The source and location of the error.
}

func (se SourceError) Error() string {
	return se.message
}

func (se SourceError) SourceAndLocation() SourceAndLocation {
	return se.sal
}

// SourceErrorf returns a new SourceError for the given location and message.
func SourceErrorf(sal SourceAndLocation, msg string, args ...interface{}) SourceError {
	return SourceError{
		message: fmt.Sprintf(msg, args...),
		sal:     sal,
	}
}

// NewSourceError returns a new SourceError for the given location and message.
func NewSourceError(sal SourceAndLocation, msg string) SourceError {
	return SourceError{
		message: msg,
		sal:     sal,
	}
}
