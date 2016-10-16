// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"github.com/serulian/compiler/generator/escommon/esbuilder"

	"github.com/cevaris/ordered_map"
)

// generatedSourceResult is the result of generating some expression source.
type generatedSourceResult struct {
	// Source is the generated source code.
	Source esbuilder.SourceBuilder

	// IsPromise indicates whether the generated source is a promise.
	IsPromise bool
}

// generatedInitMap defines an ordered map for generated source results, that also
// tracks whether any are promising.
type generatedInitMap struct {
	promising  bool
	orderedMap *ordered_map.OrderedMap
}

func newGeneratedInitMap() *generatedInitMap {
	return &generatedInitMap{
		promising:  false,
		orderedMap: ordered_map.NewOrderedMap(),
	}
}

// CombineWith combines the other map with this map.
func (gim *generatedInitMap) CombineWith(other *generatedInitMap) *generatedInitMap {
	if other.orderedMap.Len() == 0 {
		return gim
	}

	combinedMap := gim.orderedMap

	iter := other.orderedMap.IterFunc()
	for v, ok := iter(); ok; v, ok = iter() {
		combinedMap.Set(v.Key, v.Value)
	}

	return &generatedInitMap{
		promising:  gim.promising || other.promising,
		orderedMap: combinedMap,
	}
}

// Iter returns an iterator over the key/value pairs in this map.
func (gim *generatedInitMap) Iter() <-chan *ordered_map.KVPair {
	return gim.orderedMap.UnsafeIter()
}

// Promising returns whether any of the generated source items in the map are promising.
func (gim *generatedInitMap) Promising() bool {
	return gim.promising
}

// Set sets the given key to the given value. Note that the value *must* be a generatedSourceResult,
// but this function takes in an interface{} to match the interface of the normal OrderedMap.
func (gim *generatedInitMap) Set(key interface{}, value interface{}) {
	result := value.(generatedSourceResult)
	if result.IsPromise {
		gim.promising = true
		wrappedSource := esbuilder.Template("wrappedinit", `
			(init.push({{ emit . }}))
		`, result.Source)

		gim.orderedMap.Set(key, wrappedSource)
	} else {
		gim.orderedMap.Set(key, result.Source)
	}
}
