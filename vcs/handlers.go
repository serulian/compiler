// Copyright 2015 The Serulian Authors. All rights reserved.
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

// VCSKind identifies the supported kinds of VCS.
type VCSKind int

const (
	VCSKindUnknown VCSKind = iota // an unknown kind of VCS
	VCSKindGit                    // Git

	VCSKindFakeGit // Fake git for testing only
)

// vcsCheckoutFn is a function for performing a full checkout.
type vcsCheckoutFn func(path vcsPackagePath, discovery VCSUrlInformation, checkoutDir string) (string, error)

// vcsDetectFn is a function for detecting whether this handler matches the given checkout directory.
type vcsDetectFn func(checkoutDir string) bool

// vcsHasChangesFn is a function for detecting whether a directory has uncommitted code changes.
type vcsHasChangesFn func(checkoutDir string) bool

// vcsUpdateFn is a function for performing a pull/update of a checked out directory.
type vcsUpdateFn func(checkoutDir string) error

// vcsInspectFn is a function for inspecting a checked out directory.
type vcsInspectFn func(checkoutDir string) (string, error)

// vcsListTagsFn is a function for listing all tags/version of a checked out directory.
type vcsListTagsFn func(checkoutDir string) ([]string, error)

// vcsHandler represents the defined handler information for a specific kind of VCS.
type vcsHandler struct {
	kind     VCSKind         // The kind of the VCS being handled.
	checkout vcsCheckoutFn   // Function to checkout a package.
	detect   vcsDetectFn     // Function to detect if this handler matches a package.
	update   vcsUpdateFn     // Function to update a checkout.
	inspect  vcsInspectFn    // Function to inspect a checkout for its commit SHA.
	check    vcsHasChangesFn // Function to detect for code changes.
	listTags vcsListTagsFn   // Function to list tags.
}

// runCommand uses exec to run an external command.
func runCommand(runDirectory string, command string, args ...string) (error, string, string) {
	log.Printf("Running command %s %v", command, args)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Dir = runDirectory
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	log.Printf("Result: %v | %s | %s", err, stdout.String(), stderr.String())
	return err, stdout.String(), stderr.String()
}

// vcsById holds a map from string ID for the VCS to its associated vcsHandler struct.
var vcsById = map[string]vcsHandler{
	"git": vcsHandler{
		kind: VCSKindGit,

		checkout: func(vcsPath vcsPackagePath, discovery VCSUrlInformation, checkoutDir string) (string, error) {
			// Make the new directory.
			log.Printf("Making GIT package path: %s", checkoutDir)
			err := os.MkdirAll(checkoutDir, 0744)
			if err != nil {
				return "", err
			}

			// Run the clone to checkout the package.
			log.Printf("Clone git repository: %s", discovery.DownloadPath)
			if err, _, errStr := runCommand(checkoutDir, "git", "clone", discovery.DownloadPath, "."); err != nil {
				return "", fmt.Errorf("Error cloning git package %s: %s", discovery.DownloadPath, errStr)
			}

			// Checkout any submodules.
			log.Printf("Clone git repository submodules: %s", discovery.DownloadPath)

			// Note: No error check here, as submodule returns 1 on none (WHY?!)
			runCommand(checkoutDir, "git", "submodule", "update", "--init", "--recursive")

			// Switch to the tag or branch if necessary.
			switch {
			case vcsPath.branchOrCommit != "":
				log.Printf("Switch to branch on git repository: %s => %s", discovery.DownloadPath, vcsPath.branchOrCommit)
				if err, _, errStr := runCommand(checkoutDir, "git", "checkout", vcsPath.branchOrCommit); err != nil {
					return "", fmt.Errorf("Error changing branch of git package %s: %s", discovery.DownloadPath, errStr)
				}

			case vcsPath.tag != "":
				log.Printf("Switch to tag on git repository: %s => %s", discovery.DownloadPath, vcsPath.tag)
				if err, _, errStr := runCommand(checkoutDir, "git", "checkout", "tags/"+vcsPath.tag); err != nil {
					return "", fmt.Errorf("Error changing tag of git package %s: %s", discovery.DownloadPath, errStr)
				}
			}

			return "", nil
		},

		inspect: func(checkoutDir string) (string, error) {
			var out bytes.Buffer

			cmd := exec.Command("git", "rev-parse", "HEAD")
			cmd.Dir = checkoutDir
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				return "", err
			}

			trimmed := strings.TrimSpace(out.String())
			return trimmed[0:7], nil
		},

		detect: func(checkoutDir string) bool {
			gitDirectory := path.Join(checkoutDir, ".git")
			_, err := os.Stat(gitDirectory)
			return !os.IsNotExist(err)
		},

		update: func(checkoutDir string) error {
			if err, _, errStr := runCommand(checkoutDir, "git", "pull"); err != nil {
				return fmt.Errorf("Error updating git package %s: %s", checkoutDir, errStr)
			}

			return nil
		},

		check: func(checkoutDir string) bool {
			var out bytes.Buffer

			statusCmd := exec.Command("git", "status", "--porcelain")
			statusCmd.Dir = checkoutDir
			statusCmd.Stdout = &out
			statusErr := statusCmd.Run()
			if statusErr != nil {
				return true
			}

			return len(out.String()) > 0
		},

		listTags: func(checkoutDir string) ([]string, error) {
			err, output, errStr := runCommand(checkoutDir, "git", "tag")
			if err != nil {
				return []string{}, fmt.Errorf("Could not list tags: %v", errStr)
			}

			return strings.Split(output, "\n"), nil
		},
	},

	"___fake__git___": vcsHandler{
		kind: VCSKindFakeGit,

		checkout: func(vcsPath vcsPackagePath, discovery VCSUrlInformation, checkoutDir string) (string, error) {
			panic("Fake git!")
		},

		inspect: func(checkoutDir string) (string, error) {
			panic("Fake git!")
		},

		detect: func(checkoutDir string) bool {
			gitDirectory := path.Join(checkoutDir, "_dot_git")
			_, err := os.Stat(gitDirectory)
			return !os.IsNotExist(err)
		},

		update: func(checkoutDir string) error {
			panic("Fake git!")
		},

		check: func(checkoutDir string) bool {
			panic("Fake git!")
		},

		listTags: func(checkoutDir string) ([]string, error) {
			panic("Fake git!")
		},
	},
}

var vcsByKind = map[VCSKind]vcsHandler{
	VCSKindGit: vcsById["git"],

	VCSKindFakeGit: vcsById["___fake__git___"],
}

// detectHandler attempts to detect which VCS handler is applicable to the given checkout
// directory.
func detectHandler(checkoutDir string) (vcsHandler, bool) {
	for k := range vcsByKind {
		handler := vcsByKind[k]
		if handler.detect(checkoutDir) {
			return handler, true
		}
	}

	return vcsHandler{}, false
}
