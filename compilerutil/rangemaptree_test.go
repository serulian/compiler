// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type rangeMapTreeTestRange struct {
	name          string
	startPosition int
	endPosition   int
}

type rangeMapTreeTestEntry struct {
	key    string
	ranges []rangeMapTreeTestRange
}

func assertGetNil(t *testing.T, rangeMap *RangeMapTree, key string, position int) bool {
	result := rangeMap.Get(key, IntRange{position, position})
	return assert.Nil(t, result, "Expected nil result for position %s under key %s", position, key)
}

func assertGet(t *testing.T, rangeMap *RangeMapTree, key string, position int, expectedName string) bool {
	result := rangeMap.Get(key, IntRange{position, position})
	if !assert.NotNil(t, result, "Expected non-nil result for position %s under key %s", position, key) {
		return false
	}

	return assert.Equal(t, expectedName, result.(string), "Result mismatch for position %s under key %s", position, key)
}

func TestRangeMapping(t *testing.T) {
	entries := []rangeMapTreeTestEntry{
		rangeMapTreeTestEntry{
			key: "firstkey",
			ranges: []rangeMapTreeTestRange{
				rangeMapTreeTestRange{"foo", 10, 20},
				rangeMapTreeTestRange{"bar", 21, 50},
				rangeMapTreeTestRange{"baz", 100, 200},
			},
		},

		rangeMapTreeTestEntry{
			key: "secondkey",
			ranges: []rangeMapTreeTestRange{
				rangeMapTreeTestRange{"foo", 10, 25},
				rangeMapTreeTestRange{"bar", 28, 50},
				rangeMapTreeTestRange{"baz", 80, 200},
			},
		},

		rangeMapTreeTestEntry{
			key: "thirdkey",
			ranges: []rangeMapTreeTestRange{
				rangeMapTreeTestRange{"bar", 11, 20},
				rangeMapTreeTestRange{"foo", 10, 25},
			},
		},

		rangeMapTreeTestEntry{
			key: "fourthkey",
			ranges: []rangeMapTreeTestRange{
				rangeMapTreeTestRange{"foo", 10, 25},
				rangeMapTreeTestRange{"bar", 11, 20},
			},
		},
	}

	calculator := func(key string, requested IntRange) (IntRange, interface{}) {
		for _, entry := range entries {
			if entry.key == key {
				for _, current := range entry.ranges {
					if requested.StartPosition >= current.startPosition && requested.EndPosition <= current.endPosition {
						return IntRange{current.startPosition, current.endPosition}, current.name
					}
				}

				return requested, nil
			}
		}

		panic("Unknown key")
	}

	rangeMap := NewRangeMapTree(calculator)

	// Check each twice to ensuring caching is working.
	assertGet(t, rangeMap, "firstkey", 10, "foo")
	assertGet(t, rangeMap, "firstkey", 10, "foo")

	assertGet(t, rangeMap, "firstkey", 15, "foo")
	assertGet(t, rangeMap, "firstkey", 15, "foo")

	assertGet(t, rangeMap, "firstkey", 20, "foo")

	assertGet(t, rangeMap, "firstkey", 21, "bar")
	assertGet(t, rangeMap, "firstkey", 50, "bar")
	assertGetNil(t, rangeMap, "firstkey", 51)

	assertGet(t, rangeMap, "secondkey", 10, "foo")
	assertGet(t, rangeMap, "secondkey", 10, "foo")

	assertGetNil(t, rangeMap, "secondkey", 26)
	assertGetNil(t, rangeMap, "secondkey", 26)

	assertGet(t, rangeMap, "secondkey", 28, "bar")
	assertGet(t, rangeMap, "secondkey", 50, "bar")

	assertGetNil(t, rangeMap, "secondkey", 51)
	assertGetNil(t, rangeMap, "secondkey", 51)

	assertGet(t, rangeMap, "thirdkey", 12, "bar")
	assertGet(t, rangeMap, "thirdkey", 10, "foo")

	// Since `foo` is written first, we will never find `bar`.
	assertGet(t, rangeMap, "fourthkey", 12, "foo")
	assertGet(t, rangeMap, "fourthkey", 10, "foo")
}
