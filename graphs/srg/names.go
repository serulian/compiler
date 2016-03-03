// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const ASYNC_SUFFIX = "Async"

// isExportedName returns whether the given name is exported by the module to
// other modules.
func isExportedName(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

// isAsyncFunction returns whether the given name defines an async function.
func isAsyncFunction(name string) bool {
	return len(name) > len(ASYNC_SUFFIX) && strings.HasSuffix(name, ASYNC_SUFFIX)
}
