// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sourcemap

import "sort"

// orderedStringSet maintains a set of strings that preserves insertion order
// and allows fast retrieval of the index of an item inside.
type orderedStringSet struct {
	// items is a slice containing the items in their insertion order.
	items []string

	// indexMap is a map from the item to its index, if any.
	indexMap map[string]int
}

// newOrderedStringSet creates a new ordered string set.
func newOrderedStringSet() *orderedStringSet {
	return &orderedStringSet{
		items:    make([]string, 0, 25),
		indexMap: map[string]int{},
	}
}

// Add adds the given item to the set if and only if not already present.
func (oss *orderedStringSet) Add(item string) {
	if oss.IndexOf(item) >= 0 {
		// Item already exists.
		return
	}

	oss.indexMap[item] = len(oss.items)
	oss.items = append(oss.items, item)
}

// OrderedItems returns the items in insertion (or, if Sort was called, sorted) order.
func (oss *orderedStringSet) OrderedItems() []string {
	return oss.items
}

// IndexOf returns the index of the item or -1 if none.
func (oss *orderedStringSet) IndexOf(item string) int {
	if index, ok := oss.indexMap[item]; ok {
		return index
	}

	return -1
}

// Sort sorts the set.
func (oss *orderedStringSet) Sort() {
	sort.Strings(oss.items)

	// Rebuild the index map.
	for index, item := range oss.items {
		oss.indexMap[item] = index
	}
}
