// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type semvarUpdateTest struct {
	currentVersion    string
	availableVersions []string
	option            UpdateVersionOption
	resultVersion     string
	status            bool
}

var semvarUpdateTests = []semvarUpdateTest{
	// Success tests.
	semvarUpdateTest{"1.0.0", []string{}, UpdateVersionMinor, "1.0.0", true},
	semvarUpdateTest{"1.0.0", []string{"1.0.0"}, UpdateVersionMinor, "1.0.0", true},
	semvarUpdateTest{"1.0.0", []string{"1.1.0"}, UpdateVersionMinor, "1.1.0", true},
	semvarUpdateTest{"1.0.0", []string{"1.1.0pre"}, UpdateVersionMinor, "1.0.0", true},

	semvarUpdateTest{"1.0.0", []string{"1.1.0", "2.0.0"}, UpdateVersionMinor, "1.1.0", true},
	semvarUpdateTest{"1.0.0", []string{"2.1.0", "1.1.0"}, UpdateVersionMinor, "1.1.0", true},
	semvarUpdateTest{"1.0.0", []string{"1.1.0", "2.0.0"}, UpdateVersionMajor, "2.0.0", true},
	semvarUpdateTest{"1.0.0", []string{"3.1.4", "1.1.0", "2.0.0"}, UpdateVersionMajor, "3.1.4", true},

	// Failure tests.
	semvarUpdateTest{"", []string{}, UpdateVersionMinor, "", false},
	semvarUpdateTest{"a.b.c", []string{}, UpdateVersionMinor, "", false},
	semvarUpdateTest{"1.03.0", []string{}, UpdateVersionMinor, "", false},
	semvarUpdateTest{"v1.03.0", []string{}, UpdateVersionMinor, "", false},
}

func TestSemanticUpdateVersion(t *testing.T) {
	for _, test := range semvarUpdateTests {
		resultVersion, err := SemanticUpdateVersion(test.currentVersion, test.availableVersions, test.option)
		if !assert.Equal(t, test.status, err == nil, "Mismatch in status for test `%s`", test.currentVersion) {
			continue
		}

		if err != nil {
			continue
		}

		if !assert.Equal(t, test.resultVersion, resultVersion, "Mismatch in result version for test `%s`", test.currentVersion) {
			continue
		}
	}
}
