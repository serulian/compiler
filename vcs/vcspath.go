// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vcs

import (
	"fmt"
	"path"
	"regexp"
)

// vcsPackagePathRegex is a regular expression for parsing VCS paths.
//
// Examples:
//   github.com/some/project
//   github.com/some/project:somebranchorcommit
//   github.com/some/project@sometag
//   github.com/some/project//somesubdir@sometag
const vcsPackagePathRegex = "^(([a-zA-Z0-9\\._-]+)(/([a-zA-Z0-9\\._-])+)*)(//([^@:]+))?((@|:)([\\.a-zA-Z0-9_-]+))?$"

// vcsPackagePath holds information about a package located in a VCS.
type vcsPackagePath struct {
	url            string // The URL of the VCS package before discovery
	branchOrCommit string // The branch/commit of the package. Cannot be set with 'tag'.
	tag            string // The tag of the package. Cannot be set with 'branchOrCommit'.
	subpackage     string // The path of the subpackage, if any.
}

// ParseVCSPath parses a path/url to a VCS package into its components.
func ParseVCSPath(vcsPath string) (vcsPackagePath, error) {
	regex := regexp.MustCompile(vcsPackagePathRegex)
	matches := regex.FindStringSubmatch(vcsPath)
	if matches == nil {
		return vcsPackagePath{}, fmt.Errorf("Invalid VCS package path: %s", vcsPath)
	}

	cleaned := path.Clean(matches[1])
	if cleaned != matches[1] {
		return vcsPackagePath{}, fmt.Errorf("Invalid VCS package path: %s", vcsPath)
	}

	var branchOrComment string
	var tag string

	if matches[8] == ":" {
		branchOrComment = matches[9]
	} else if matches[8] == "@" {
		tag = matches[9]
	}

	return vcsPackagePath{cleaned, branchOrComment, tag, matches[6]}, nil
}

// URL returns the URL of the VCS package.
func (pp vcsPackagePath) URL() string {
	return pp.url
}

// Tag returns the tag of the VCS package.
func (pp vcsPackagePath) Tag() string {
	return pp.tag
}

// BranchOrCommit returns the branch or commit of the VCS package.
func (pp vcsPackagePath) BranchOrCommit() string {
	return pp.branchOrCommit
}

// WithCommit returns the VCS package path with the given commit.
func (pp vcsPackagePath) WithCommit(commitSha string) vcsPackagePath {
	return vcsPackagePath{
		url:            pp.url,
		branchOrCommit: commitSha,
		tag:            "",
		subpackage:     pp.subpackage,
	}
}

// WithTag returns the VCS package path with the given tag.
func (pp vcsPackagePath) WithTag(tag string) vcsPackagePath {
	return vcsPackagePath{
		url:            pp.url,
		branchOrCommit: "",
		tag:            tag,
		subpackage:     pp.subpackage,
	}
}

// AsGeneric returns the VCS import without any commit, branch or tag.
func (pp vcsPackagePath) AsGeneric() vcsPackagePath {
	return vcsPackagePath{
		url:            pp.url,
		branchOrCommit: "",
		tag:            "",
		subpackage:     pp.subpackage,
	}
}

// String returns the string representation of this path.
func (pp vcsPackagePath) String() string {
	var strRep = pp.url

	if pp.subpackage != "" {
		strRep = strRep + "//" + pp.subpackage
	}

	switch {
	case pp.branchOrCommit != "":
		return strRep + ":" + pp.branchOrCommit

	case pp.tag != "":
		return strRep + "@" + pp.tag

	default:
		return strRep
	}
}

// isHEAD returns whether this VCS package is pointed to HEAD of a branch or a specific commit.
func (pp vcsPackagePath) isHEAD() bool {
	return pp.tag == ""
}

// isRepoOnlyReference returns true if this VCS package does not specify its tag, branch or commit.
func (pp vcsPackagePath) isRepoOnlyReference() bool {
	return pp.tag == "" && pp.branchOrCommit == ""
}

// cacheDirectory returns a relative directory path at which the VCS package should be placed
// for caching.
func (pp vcsPackagePath) cacheDirectory() string {
	var suffix string
	switch {
	case pp.branchOrCommit != "":
		suffix = "branch/" + pp.branchOrCommit

	case pp.tag != "":
		suffix = "tag/" + pp.tag

	default:
		suffix = "HEAD"
	}

	return path.Join(pp.url, suffix)
}
