// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sourcemap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type vlqTest struct {
	encoded        string
	decoded        []int
	skipEncodeTest bool
}

var vlqTests = []vlqTest{
	vlqTest{"A", []int{0}, false},
	vlqTest{"B", []int{0}, true}, // Skip encoding since 0 = 'A'
	vlqTest{"C", []int{1}, false},
	vlqTest{"f", []int{-15}, false},

	vlqTest{"gA", []int{0}, true}, // Skip encoding since 0 = 'A'
	vlqTest{"gB", []int{16}, false},
	vlqTest{"RgB", []int{-8, 16}, false},
	vlqTest{"EAEE", []int{2, 0, 2, 2}, false},
	vlqTest{"OAAQ", []int{7, 0, 0, 8}, false},
	vlqTest{"AAgBC", []int{0, 0, 16, 1}, false},
	vlqTest{"SACjBD", []int{9, 0, 1, -17, -1}, false},
	vlqTest{"kBAAhBA", []int{18, 0, 0, -16, 0}, false},
	vlqTest{"EAEEga", []int{2, 0, 2, 2, 416}, false},
}

func TestVLQ(t *testing.T) {
	// Test decoding and encoding.
	for _, test := range vlqTests {
		values, ok := vlqDecode(test.encoded)
		if !assert.True(t, ok, "Expected proper decode for %s", test.encoded) {
			continue
		}

		if !assert.Equal(t, test.decoded, values, "Decode mismatch for %s", test.encoded) {
			continue
		}

		if test.skipEncodeTest {
			continue
		}

		encoded := vlqEncode(test.decoded)

		if !assert.Equal(t, test.encoded, encoded, "Encode mismatch for %s", test.encoded) {
			continue
		}
	}
}
