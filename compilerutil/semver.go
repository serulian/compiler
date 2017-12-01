// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

import (
	"github.com/blang/semver"
)

// UpdateVersionOption defines the options for the SemanticUpdateVersion method.
type UpdateVersionOption string

const (
	// UpdateVersionMinor indicates that only minor version updates are allowed.
	UpdateVersionMinor UpdateVersionOption = "minor"

	// UpdateVersionMajor indicates the both minor and major version updates are allowed.
	UpdateVersionMajor UpdateVersionOption = "major"
)

// NextVersionOption defines the options for the NextSemanticVersion method.
type NextVersionOption string

const (
	// NextMajorVersion indicates that the major version field should be incremented.
	NextMajorVersion NextVersionOption = "major"

	// NextMinorVersion indicates that the minor version field should be incremented.
	NextMinorVersion NextVersionOption = "minor"

	// NextPatchVersion indicates that the patch version field should be incremented.
	NextPatchVersion NextVersionOption = "patch"
)

// NextSemanticVersion returns the next logic semantic version from the current version, given the
// update option.
func NextSemanticVersion(currentVersionString string, option NextVersionOption) (string, error) {
	version, err := semver.ParseTolerant(currentVersionString)
	if err != nil {
		return "", err
	}

	// Clear the Pre and Build fields.
	version.Pre = []semver.PRVersion{}
	version.Build = []string{}

	switch option {
	case NextMajorVersion:
		version.Major++
		version.Minor = 0
		version.Patch = 0

	case NextMinorVersion:
		version.Minor++
		version.Patch = 0

	case NextPatchVersion:
		version.Patch++

	default:
		panic("Unknown NextVersionOption")
	}

	return version.String(), nil
}

// LatestSemanticVersion returns the latest non-patch semantic version found in the list of available
// versions, if any.
func LatestSemanticVersion(availableVersions []string) (string, bool) {
	latestVersionString := ""
	for _, possibleVersion := range availableVersions {
		// Skip empty version strings.
		if len(possibleVersion) == 0 {
			continue
		}

		// Skip possibleVersions that don't parse, as well as pre-release versions
		// (since they are, by definition, not considered stable).
		parsed, err := semver.ParseTolerant(possibleVersion)
		if err != nil || len(parsed.Pre) > 0 {
			continue
		}

		latestVersionString = possibleVersion
	}

	return latestVersionString, len(latestVersionString) > 0
}

// SemanticUpdateVersion finds the semver in the available versions list that is newer than the current
// version.
func SemanticUpdateVersion(currentVersionString string, availableVersions []string, option UpdateVersionOption) (string, error) {
	currentVersion, err := semver.ParseTolerant(currentVersionString)
	if err != nil {
		return "", err
	}

	latestVersionString := currentVersionString
	for _, possibleVersion := range availableVersions {
		// Skip empty version strings.
		if len(possibleVersion) == 0 {
			continue
		}

		// Skip possibleVersions that don't parse, as well as pre-release versions
		// (since they are, by definition, not considered stable).
		parsed, err := semver.ParseTolerant(possibleVersion)
		if err != nil || len(parsed.Pre) > 0 {
			continue
		}

		// Find the latest stable version.
		if parsed.GT(currentVersion) {
			if option == UpdateVersionMajor || parsed.Major == currentVersion.Major {
				currentVersion = parsed
				latestVersionString = possibleVersion
			}
		}
	}

	return latestVersionString, nil
}
