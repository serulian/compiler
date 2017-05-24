// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilercommon

import (
	"fmt"
)

// SourceError represents an error produced by a source file at a specific range.
type SourceError struct {
	message     string
	sourceRange SourceRange
}

func (se SourceError) Error() string {
	return se.message
}

func (se SourceError) String() string {
	return se.message
}

// SourceRange returns the range of this error.
func (se SourceError) SourceRange() SourceRange {
	return se.sourceRange
}

// SourceErrorf returns a new SourceError for the given range and message.
func SourceErrorf(sourceRange SourceRange, msg string, args ...interface{}) SourceError {
	return SourceError{
		message:     fmt.Sprintf(msg, args...),
		sourceRange: sourceRange,
	}
}

// NewSourceError returns a new SourceError for the given range and message.
func NewSourceError(sourceRange SourceRange, msg string) SourceError {
	return SourceError{
		message:     msg,
		sourceRange: sourceRange,
	}
}
