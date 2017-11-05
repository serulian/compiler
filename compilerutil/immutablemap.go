// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

import (
	"hash/fnv"

	hamt "github.com/raviqqe/hamt.go"
)

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
	return copyImmutableMap{
		internalMap: map[string]interface{}{},
	}
}

// hamtMapThreshold defines the threshold at which we switch from a copy map
// to an hamt map.
const hamtMapThreshold = 10

// copyImmutableMap is an immutable map that uses the internal Go map
// and fully copies it everytime a Set operation is called. This is quite
// fast for maps of size <= 10, after which it starts to lose some performance
// and really gets bad at sizes >= 50.
type copyImmutableMap struct {
	internalMap map[string]interface{}
}

// hamtImmutableMap is an immutable map thatuses the hamt package's
// Map implementation. This has better performance for maps of size >= 50.
type hamtImmutableMap struct {
	internalMap hamt.Map
}

func (i copyImmutableMap) Get(key string) (interface{}, bool) {
	value, ok := i.internalMap[key]
	return value, ok
}

func (i copyImmutableMap) Set(key string, value interface{}) ImmutableMap {
	// If we've reached the threshold, switch from a copy map to the hamt map.
	length := len(i.internalMap)
	if length >= hamtMapThreshold {
		var hamtMap ImmutableMap = hamtImmutableMap{hamt.NewMap()}
		for existingKey, value := range i.internalMap {
			hamtMap = hamtMap.Set(existingKey, value)
		}
		hamtMap = hamtMap.Set(key, value)
		return hamtMap
	}

	// Otherwise, create a new map, copy over all the existing elements, and
	// set the new key.
	newMap := make(map[string]interface{}, length)
	for existingKey, value := range i.internalMap {
		newMap[existingKey] = value
	}
	newMap[key] = value
	return copyImmutableMap{newMap}
}

func (i hamtImmutableMap) Get(key string) (interface{}, bool) {
	value := i.internalMap.Find(hamtString(key))
	return value, value != nil
}

func (i hamtImmutableMap) Set(key string, value interface{}) ImmutableMap {
	return hamtImmutableMap{i.internalMap.Insert(hamtString(key), value)}
}

type hamtString string

func (h hamtString) Equal(other hamt.Entry) bool {
	return string(h) == string(other.(hamtString))
}

func (h hamtString) Hash() uint32 {
	hsh := fnv.New32a()
	hsh.Write([]byte(string(h)))
	return hsh.Sum32()
}
