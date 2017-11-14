// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vcs

import (
	"os"
	"path"
)

type fakeGitVcs struct{}

func (fgv fakeGitVcs) Kind() string {
	return "__fake_git__"
}

func (fgv fakeGitVcs) Detect(checkoutDir string) bool {
	gitDirectory := path.Join(checkoutDir, "_dot_git")
	_, err := os.Stat(gitDirectory)
	return !os.IsNotExist(err)
}

func (fgv fakeGitVcs) Checkout(path vcsPackagePath, downloadPath string, checkoutDir string) error {
	panic("Fake git!")
}

func (fgv fakeGitVcs) HasLocalChanges(checkoutDir string) bool {
	panic("Fake git!")
}

func (fgv fakeGitVcs) Update(checkoutDir string) error {
	panic("Fake git!")
}

func (fgv fakeGitVcs) Inspect(checkoutDir string) (string, error) {
	panic("Fake git!")
}

func (fgv fakeGitVcs) ListTags(checkoutDir string) ([]string, error) {
	panic("Fake git!")
}

func (fgv fakeGitVcs) IsDetached(checkoutDir string) (bool, error) {
	panic("Fake git!")
}
