// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vcs

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type gitVcs struct{}

func (gv gitVcs) Kind() string {
	return "git"
}

func (gv gitVcs) Checkout(vcsPath vcsPackagePath, downloadPath string, checkoutDir string) error {
	lock := modificationLockMap.GetLock(checkoutDir)
	lock.Lock()
	defer lock.Unlock()

	// Make the new directory.
	log.Printf("Making GIT package path: %s", checkoutDir)
	err := os.MkdirAll(checkoutDir, 0744)
	if err != nil {
		return err
	}

	// Run the clone to checkout the package.
	log.Printf("Clone git repository: %s", downloadPath)
	if _, errStr, err := runCommand(checkoutDir, "git", "clone", downloadPath, "."); err != nil {
		return fmt.Errorf("Error cloning git package %s: %s", downloadPath, errStr)
	}

	// Checkout any submodules.
	log.Printf("Clone git repository submodules: %s", downloadPath)

	// Note: No error check here, as submodule returns 1 on none (WHY?!)
	runCommand(checkoutDir, "git", "submodule", "update", "--init", "--recursive")

	// Switch to the tag or branch if necessary.
	switch {
	case vcsPath.branchOrCommit != "":
		log.Printf("Switch to branch on git repository: %s => %s", downloadPath, vcsPath.branchOrCommit)
		if _, errStr, err := runCommand(checkoutDir, "git", "checkout", vcsPath.branchOrCommit); err != nil {
			return fmt.Errorf("Error changing branch of git package %s: %s", downloadPath, errStr)
		}

	case vcsPath.tag != "":
		log.Printf("Switch to tag on git repository: %s => %s", downloadPath, vcsPath.tag)
		if _, errStr, err := runCommand(checkoutDir, "git", "checkout", "tags/"+vcsPath.tag); err != nil {
			return fmt.Errorf("Error changing tag of git package %s: %s", downloadPath, errStr)
		}
	}

	return nil
}

func (gv gitVcs) Inspect(checkoutDir string) (string, error) {
	var out bytes.Buffer

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = checkoutDir
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(out.String(), "fatal:") {
		return "", fmt.Errorf("Invalid repository")
	}

	trimmed := strings.TrimSpace(out.String())
	return trimmed[0:7], nil
}

func (gv gitVcs) Detect(checkoutDir string) bool {
	gitDirectory := path.Join(checkoutDir, ".git")
	_, err := os.Stat(gitDirectory)
	return !os.IsNotExist(err)
}

func (gv gitVcs) Update(checkoutDir string) error {
	lock := modificationLockMap.GetLock(checkoutDir)
	lock.Lock()
	defer lock.Unlock()

	if _, errStr, err := runCommand(checkoutDir, "git", "pull"); err != nil {
		return fmt.Errorf("Error updating git package %s: %s", checkoutDir, errStr)
	}

	return nil
}

func (gv gitVcs) HasLocalChanges(checkoutDir string) bool {
	var out bytes.Buffer

	statusCmd := exec.Command("git", "status", "--porcelain")
	statusCmd.Dir = checkoutDir
	statusCmd.Stdout = &out
	statusErr := statusCmd.Run()
	if statusErr != nil {
		return true
	}

	return len(out.String()) > 0
}

func (gv gitVcs) ListTags(checkoutDir string) ([]string, error) {
	output, errStr, err := runCommand(checkoutDir, "git", "ls-remote", "--tags", "--refs", "--q")
	if err != nil {
		return []string{}, fmt.Errorf("Could not list tags: %v", errStr)
	}

	lines := strings.Split(output, "\n")
	tags := make([]string, 0, len(lines))
	for _, line := range lines {
		pieces := strings.Split(line, "\t")
		if len(pieces) != 2 {
			continue
		}

		ref := pieces[1]
		if !strings.HasPrefix(ref, "refs/tags/") {
			continue
		}

		tags = append(tags, ref[len("refs/tags/"):])
	}

	return tags, nil
}

func (gv gitVcs) IsDetached(checkoutDir string) (bool, error) {
	_, errStr, err := runCommand(checkoutDir, "git", "symbolic-ref", "HEAD")
	if err != nil {
		if strings.Contains(errStr, "not a symbolic ref") {
			return true, nil
		}

		return false, fmt.Errorf("Could check repo status: %v", errStr)
	}

	return false, nil
}

func (gv gitVcs) GetPackagePath(checkoutDir string) (vcsPackagePath, error) {
	// Find the origin URL.
	originUrlStr, _, err := runCommand(checkoutDir, "git", "remote", "get-url", "origin")
	if err != nil {
		return vcsPackagePath{}, err
	}

	// Check for a branch.
	branchStr, _, err := runCommand(checkoutDir, "git", "symbolic-ref", "-q", "--short", "HEAD")
	if err != nil {
		return vcsPackagePath{}, err
	}

	currentBranch := strings.TrimSpace(branchStr)
	if len(currentBranch) > 0 {
		return vcsPackagePath{
			url:            strings.TrimSpace(originUrlStr),
			branchOrCommit: branchStr,
			tag:            "",
			subpackage:     "",
		}, nil
	}

	// Check for a tag.
	tagStr, _, err := runCommand(checkoutDir, "git", "describe", "--tags", "--exact-match")
	if err != nil {
		return vcsPackagePath{}, err
	}

	currentTag := strings.TrimSpace(tagStr)
	if len(currentTag) > 0 {
		return vcsPackagePath{
			url:            strings.TrimSpace(originUrlStr),
			branchOrCommit: "",
			tag:            currentTag,
			subpackage:     "",
		}, nil
	}

	// Use the inspect result.
	currentSHA, err := gv.Inspect(checkoutDir)
	if err != nil {
		return vcsPackagePath{}, err
	}

	return vcsPackagePath{
		url:            strings.TrimSpace(originUrlStr),
		branchOrCommit: currentSHA,
		tag:            "",
		subpackage:     "",
	}, nil
}
