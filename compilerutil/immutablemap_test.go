// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertNotFound(t *testing.T, immutableMap ImmutableMap, key string) bool {
	_, found := immutableMap.Get(key)
	return assert.False(t, found, "Did not expect key %s", key)
}

func assertFound(t *testing.T, immutableMap ImmutableMap, key string, value interface{}) bool {
	foundValue, found := immutableMap.Get(key)
	return assert.True(t, found, "Expected key %s", key) &&
		assert.Equal(t, foundValue, value, "Mismatch in value for key %s", key)
}

func TestImmutableMap(t *testing.T) {
	immutableMap := NewImmutableMap()

	someKeyMap := immutableMap.Set("somekey", "somevalue")
	assertNotFound(t, immutableMap, "somekey")
	assertFound(t, someKeyMap, "somekey", "somevalue")

	anotherKeyMap := immutableMap.Set("anotherkey", "anothervalue")
	assertNotFound(t, immutableMap, "anotherkey")
	assertNotFound(t, someKeyMap, "anotherkey")

	assertFound(t, anotherKeyMap, "anotherkey", "anothervalue")
	assertNotFound(t, anotherKeyMap, "somekey")

	chainedMap := someKeyMap.Set("thirdkey", "thirdvalue")
	assertNotFound(t, immutableMap, "thirdkey")
	assertNotFound(t, someKeyMap, "thirdkey")
	assertNotFound(t, anotherKeyMap, "thirdkey")

	assertFound(t, chainedMap, "somekey", "somevalue")
	assertFound(t, chainedMap, "thirdkey", "thirdvalue")

	overwriteMap := chainedMap.Set("somekey", "newvalue")
	assertFound(t, someKeyMap, "somekey", "somevalue")

	assertFound(t, overwriteMap, "somekey", "newvalue")
	assertFound(t, overwriteMap, "thirdkey", "thirdvalue")
}

func benchmarkImmutableMap(size int, b *testing.B) {
	immutableMap := NewImmutableMap()

	for n := 0; n < b.N; n++ {
		for i := 0; i < size; i++ {
			immutableMap = immutableMap.Set(strconv.Itoa(i), i)
		}

		for i := 0; i < size; i++ {
			_, found := immutableMap.Get(strconv.Itoa(i))
			assert.True(b, found)
		}

		for i := 0; i < size/2; i++ {
			_, found := immutableMap.Get("f" + strconv.Itoa(i))
			assert.False(b, found)
		}
	}
}

func BenchmarkImmutableMap1(b *testing.B)   { benchmarkImmutableMap(1, b) }
func BenchmarkImmutableMap2(b *testing.B)   { benchmarkImmutableMap(2, b) }
func BenchmarkImmutableMap3(b *testing.B)   { benchmarkImmutableMap(3, b) }
func BenchmarkImmutableMap5(b *testing.B)   { benchmarkImmutableMap(5, b) }
func BenchmarkImmutableMap10(b *testing.B)  { benchmarkImmutableMap(10, b) }
func BenchmarkImmutableMap20(b *testing.B)  { benchmarkImmutableMap(20, b) }
func BenchmarkImmutableMap50(b *testing.B)  { benchmarkImmutableMap(50, b) }
func BenchmarkImmutableMap100(b *testing.B) { benchmarkImmutableMap(100, b) }
