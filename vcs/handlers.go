// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vcs

// VCSHandler defines a handler for working with a VCS package.
type VCSHandler interface {
	// Kind returns the kind of this handler.
	Kind() string

	// Detect detects whether this handler matches the given checkout directory.
	Detect(checkoutDir string) bool

	// Checkout performs a full checkout.
	Checkout(path vcsPackagePath, downloadPath string, checkoutDir string) error

	// HasLocalChanges detects whether the directory has uncommitted code changes.
	HasLocalChanges(checkoutDir string) bool

	// Update performs a pull/update of a checked out directory.
	Update(checkoutDir string) error

	// Inspect inspects a checked out directory, returning its HEAD SHA.
	Inspect(checkoutDir string) (string, error)

	// ListTags lists all tags/version of a project in a checked out directory.
	ListTags(checkoutDir string) ([]string, error)

	// IsDetached returns whether the checked out directory is in a detached state.
	IsDetached(checkoutDir string) (bool, error)

	// GetPackagePath returns the VCS Package Path information for the checked out directory.
	GetPackagePath(checkoutDir string) (vcsPackagePath, error)
}

var vcsByKind = map[string]VCSHandler{
	"git":          gitVcs{},
	"__fake_git__": fakeGitVcs{},
}

// GetHandlerByKind returns the VCS handler for the given VCS kind, if any.
func GetHandlerByKind(kind string) (VCSHandler, bool) {
	found, ok := vcsByKind[kind]
	return found, ok
}

// DetectHandler attempts to detect which VCS handler is applicable to the given checkout
// directory.
func DetectHandler(checkoutDir string) (VCSHandler, bool) {
	for _, handler := range vcsByKind {
		if handler.Detect(checkoutDir) {
			return handler, true
		}
	}

	return nil, false
}
