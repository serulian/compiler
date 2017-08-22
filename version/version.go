// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package version defines the versioning information for the toolkit.
package version

var (
	// Version is the version of the toolkit. Will be empty when developing.
	Version = ""

	// GitSHA is the GIT SHA of this version of the toolkit. Will be empty when developing.
	GitSHA = ""
)

// CoreLibraryTagOrBranch returns the core library tag or branch to use for this version of
// the toolkit.
func CoreLibraryTagOrBranch() string {
	if Version == "" {
		return ""
	}

	return "@" + Version
}

// DescriptiveVersion returns a descriptive form of the version that will never be empty.
func DescriptiveVersion() string {
	if Version == "" {
		return "(DEVELOPMENT)"
	}

	return Version
}
