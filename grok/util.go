// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"github.com/serulian/compiler/compilercommon"
)

type sourceLocationCapable interface {
	SourceLocation() (compilercommon.SourceAndLocation, bool)
}

func getSAL(slc sourceLocationCapable) *compilercommon.SourceAndLocation {
	sal, hasSal := slc.SourceLocation()
	if hasSal {
		return &sal
	}

	return nil
}
