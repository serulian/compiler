// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate protoc --gofast_out=. sourcelocation.proto

package proto

import (
	"github.com/serulian/compiler/compilercommon"
)

func (t *SourceLocation) Name() string {
	return "SourceLocation"
}

func (t *SourceLocation) Value() string {
	bytes, err := t.Marshal()
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func (t *SourceLocation) Build(value string) interface{} {
	uerr := t.Unmarshal([]byte(value))
	if uerr != nil {
		panic(uerr)
	}

	return t
}

// AsSourceAndLocation returns the SourceAndLocation representation.
func (t *SourceLocation) AsSourceAndLocation() compilercommon.SourceAndLocation {
	return compilercommon.NewSourceAndLocation(compilercommon.InputSource(t.SourcePath), int(t.BytePosition))
}
