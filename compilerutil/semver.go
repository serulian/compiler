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
	UpdateVersionMajor = "major"
)

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
