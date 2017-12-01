// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type semverNextTest struct {
	currentVersion string
	option         NextVersionOption
	resultVersion  string
}

var semverNextTests = []semverNextTest{
	semverNextTest{"1.0.0", NextPatchVersion, "1.0.1"},
	semverNextTest{"1.0.0", NextMinorVersion, "1.1.0"},
	semverNextTest{"1.0.0", NextMajorVersion, "2.0.0"},

	semverNextTest{"1.1.2", NextPatchVersion, "1.1.3"},
	semverNextTest{"1.1.2", NextMinorVersion, "1.2.0"},
	semverNextTest{"1.1.2", NextMajorVersion, "2.0.0"},

	semverNextTest{"1.1.2-foo", NextPatchVersion, "1.1.3"},
	semverNextTest{"1.1.2-foo", NextMinorVersion, "1.2.0"},
	semverNextTest{"1.1.2-foo", NextMajorVersion, "2.0.0"},

	semverNextTest{"1.1.2+foo", NextPatchVersion, "1.1.3"},
	semverNextTest{"1.1.2+foo", NextMinorVersion, "1.2.0"},
	semverNextTest{"1.1.2+foo", NextMajorVersion, "2.0.0"},
}

func TestNextSemanticVersion(t *testing.T) {
	for _, test := range semverNextTests {
		resultVersion, err := NextSemanticVersion(test.currentVersion, test.option)
		if !assert.Nil(t, err, "Mismatch in status for next test `%s`", test.currentVersion) {
			continue
		}

		if !assert.Equal(t, test.resultVersion, resultVersion, "Mismatch in result version for next test `%s`", test.currentVersion) {
			continue
		}
	}
}

type semverUpdateTest struct {
	currentVersion    string
	availableVersions []string
	option            UpdateVersionOption
	resultVersion     string
	status            bool
}

var semverUpdateTests = []semverUpdateTest{
	// Success tests.
	semverUpdateTest{"1.0.0", []string{}, UpdateVersionMinor, "1.0.0", true},
	semverUpdateTest{"1.0.0", []string{"1.0.0"}, UpdateVersionMinor, "1.0.0", true},
	semverUpdateTest{"1.0.0", []string{"1.1.0"}, UpdateVersionMinor, "1.1.0", true},
	semverUpdateTest{"1.0.0", []string{"1.1.0pre"}, UpdateVersionMinor, "1.0.0", true},

	semverUpdateTest{"1.0.0", []string{"1.1.0", "2.0.0"}, UpdateVersionMinor, "1.1.0", true},
	semverUpdateTest{"1.0.0", []string{"2.1.0", "1.1.0"}, UpdateVersionMinor, "1.1.0", true},
	semverUpdateTest{"1.0.0", []string{"1.1.0", "2.0.0"}, UpdateVersionMajor, "2.0.0", true},
	semverUpdateTest{"1.0.0", []string{"3.1.4", "1.1.0", "2.0.0"}, UpdateVersionMajor, "3.1.4", true},

	// Failure tests.
	semverUpdateTest{"", []string{}, UpdateVersionMinor, "", false},
	semverUpdateTest{"a.b.c", []string{}, UpdateVersionMinor, "", false},
	semverUpdateTest{"1.03.0", []string{}, UpdateVersionMinor, "", false},
	semverUpdateTest{"v1.03.0", []string{}, UpdateVersionMinor, "", false},
}

func TestSemanticUpdateVersion(t *testing.T) {
	for _, test := range semverUpdateTests {
		resultVersion, err := SemanticUpdateVersion(test.currentVersion, test.availableVersions, test.option)
		if !assert.Equal(t, test.status, err == nil, "Mismatch in status for update test `%s`", test.currentVersion) {
			continue
		}

		if err != nil {
			continue
		}

		if !assert.Equal(t, test.resultVersion, resultVersion, "Mismatch in result version for update test `%s`", test.currentVersion) {
			continue
		}
	}
}
