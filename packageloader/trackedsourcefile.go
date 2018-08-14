// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import "github.com/serulian/compiler/compilercommon"

// trackedSourceFile defines a struct for tracking the versioning and contents of a source file
// encountered during package loading.
type trackedSourceFile struct {
	// sourcePath is the path of the source file encountered.
	sourcePath compilercommon.InputSource

	// sourceKind is the kind of the source file.
	sourceKind string

	// contents are the contents of the source file when it was loaded.
	contents []byte

	// revisionID is the revision ID of the source file when it was loaded, as returned by
	// GetRevisionID on the PathLoader. Typically, this is an mtime or version number. If the
	// revision ID has changed, the contents of the file are assumed to have been altered since
	// the file was read.
	revisionID int64

	// trackedPositionMapper is the source position mapper for converting rune positions <-> line+col
	// positions. Only instantiated if necessary and only instantiated for the *tracked* contents of
	// the file.
	trackedPositionMapper *positionMapperAndRevisionID

	// latestPositionMapper is the source position mapper for converting rune positions <-> line+col
	// positions. This will be reinstantiated as necessary when the current contents of the file have
	// changed.
	latestPositionMapper *positionMapperAndRevisionID
}

// positionMapperAndRevisionID holds a position mapper and the file revision to which it applies.
type positionMapperAndRevisionID struct {
	revisionID     int64
	positionMapper compilercommon.SourcePositionMapper
}
