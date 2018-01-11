// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bundle

type bundler struct {
	files map[string]BundledFile
}

// NewBundler returns a new file bundler.
func NewBundler() Bundler {
	return &bundler{
		files: map[string]BundledFile{},
	}
}

func (b *bundler) AddFile(file BundledFile) {
	b.files[file.Filename()] = file
}

func (b *bundler) Freeze(builder Builder) Bundle {
	return builder(b.files)
}
