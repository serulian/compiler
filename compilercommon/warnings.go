// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilercommon

import (
	"fmt"
)

type Warning interface {
	Warning() string // Returns the warning message.
	String() string
}

// SourceWarning represents an warning produced by a source file at a specific location.
type SourceWarning struct {
	message string            // The warning message.
	sal     SourceAndLocation // The source and location of the error.
}

func (sw *SourceWarning) Warning() string {
	return sw.message
}

func (sw *SourceWarning) String() string {
	return sw.message
}

// SourcePositionWarningf returns a new SourceWarning for the given source, byte position and message.
func SourcePositionWarningf(source InputSource, bytePosition int, msg string, args ...interface{}) *SourceWarning {
	return SourceWarningf(NewSourceAndLocation(source, bytePosition), msg, args...)
}

// SourceWarningf returns a new SourceWarning for the given location and message.
func SourceWarningf(sal SourceAndLocation, msg string, args ...interface{}) *SourceWarning {
	return &SourceWarning{
		message: fmt.Sprintf(msg, args...),
		sal:     sal,
	}
}

// NewSourceWarning returns a new SourceWarning for the given location and message.
func NewSourceWarning(sal SourceAndLocation, msg string) *SourceWarning {
	return &SourceWarning{
		message: msg,
		sal:     sal,
	}
}
