// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilercommon

import (
	"fmt"
)

// SourceWarning represents an warning produced by a source file at a specific range.
type SourceWarning struct {
	message     string
	sourceRange SourceRange
}

// Warning returns the warning message.
func (sw SourceWarning) Warning() string {
	return sw.message
}

func (sw SourceWarning) String() string {
	return sw.message
}

// SourceRange returns the range of this warning.
func (sw SourceWarning) SourceRange() SourceRange {
	return sw.sourceRange
}

// SourceWarningf returns a new SourceWarning for the given range and message.
func SourceWarningf(sourceRange SourceRange, msg string, args ...interface{}) SourceWarning {
	return SourceWarning{
		message:     fmt.Sprintf(msg, args...),
		sourceRange: sourceRange,
	}
}

// NewSourceWarning returns a new SourceWarning for the given range and message.
func NewSourceWarning(sourceRange SourceRange, msg string) SourceWarning {
	return SourceWarning{
		message:     msg,
		sourceRange: sourceRange,
	}
}
