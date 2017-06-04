// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilercommon

// CodeSummary holds the code and optional documentation of a member, parameter or other entity in the compiler.
// Used by IDE tooling for human-readable summarization.
type CodeSummary struct {
	// Documentation is the documentation string for the code summary, if any.
	Documentation string

	// Code is the code-view of the summarized entity.
	Code string

	// HasDeclaredType returns whether the code summary's code entry already contains a type.
	HasDeclaredType bool
}
