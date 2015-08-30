// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

// inputSource represents the source of parsing input (a file, etc).
type InputSource string

// bytePosition represents the byte position in a piece of code.
type bytePosition int

// sourcePosition represents a location in an input source
type sourcePosition struct {
	lineNumber     int
	columnPosition int
}
