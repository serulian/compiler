// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vcs

import (
	"bytes"
	"log"
	"os/exec"

	"github.com/serulian/compiler/compilerutil"
)

// modificationLockMap is a map for locking checkouts or updates of specific paths, to prevent multiple concurrent calls
// from being made.
var modificationLockMap = compilerutil.CreateLockMap()

// runCommand uses exec to run an external command.
func runCommand(runDirectory string, command string, args ...string) (string, string, error) {
	log.Printf("Running command %s %v", command, args)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Dir = runDirectory
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	log.Printf("Result: %v | %s | %s", err, stdout.String(), stderr.String())
	return stdout.String(), stderr.String(), err
}
