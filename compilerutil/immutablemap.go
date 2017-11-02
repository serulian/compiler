// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

// ImmutableMap defines an immutable map struct, where Set-ing a new key returns a new ImmutableMap.
type ImmutableMap interface {
	// Get returns the value found at the given key, if any.
	Get(key string) (interface{}, bool)

	// Set returns a new ImmutableMap which is a copy of this map, but with the given key set
	// to the given value.
	Set(key string, value interface{}) ImmutableMap
}

// NewImmutableMap creates a new, empty immutable map.
func NewImmutableMap() ImmutableMap {
	return immutableMap{
		internalMap: map[string]interface{}{},
	}
}

type immutableMap struct {
	internalMap map[string]interface{}
}

func (i immutableMap) Get(key string) (interface{}, bool) {
	value, ok := i.internalMap[key]
	return value, ok
}

func (i immutableMap) Set(key string, value interface{}) ImmutableMap {
	newMap := map[string]interface{}{}
	for existingKey, value := range i.internalMap {
		newMap[existingKey] = value
	}

	newMap[key] = value
	return immutableMap{newMap}
}
