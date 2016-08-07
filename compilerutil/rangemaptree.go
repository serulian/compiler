// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

import (
	"sync"

	"github.com/emirpasic/gods/trees/redblacktree"
)

// RangeMapTree defines a helper struct for cached, faster lookup of an int range to an associated
// piece of data, under a specific key (such as a source file path). Note that the ranges for the
// data must be explicitly *non-overlapping* under each key, or this data structure will fail to
// work as intended.
//
// This data structure is safe for multi-threaded access.
type RangeMapTree struct {
	//  from a particular input source to its associated mapper.
	rangeTreeMap map[string]rangeTreeEntry
	globalLock   *sync.RWMutex
	calculator   RangeMapTreeCalculator
}

// NewRangeMapTree creates a new range map tree.
func NewRangeMapTree(calculator RangeMapTreeCalculator) *RangeMapTree {
	return &RangeMapTree{
		rangeTreeMap: map[string]rangeTreeEntry{},
		globalLock:   &sync.RWMutex{},
		calculator:   calculator,
	}
}

// IntRange defines a range of integer values, inclusive.
type IntRange struct {
	StartPosition int
	EndPosition   int
}

// RangeMapTreeCalculator returns the calculated value for the given key and range. The range
// returned should encompass the entire range of the returned value. If no valid value is found
// for the position and key, the range returned should be the range requested.
type RangeMapTreeCalculator func(key string, current IntRange) (IntRange, interface{})

// Get returns the data associated with the given key and range. If the data for that range
// has not yet been calculated, the calculator method is called to do so lazily.
func (rmt *RangeMapTree) Get(key string, current IntRange) interface{} {
	rmt.globalLock.RLock()
	currentEntry, hasKey := rmt.rangeTreeMap[key]
	rmt.globalLock.RUnlock()

	if !hasKey {
		currentEntry = rangeTreeEntry{
			key:          key,
			internalTree: redblacktree.NewWith(rangeComparator),
			entryLock:    &sync.RWMutex{},
		}

		rmt.globalLock.Lock()
		rmt.rangeTreeMap[key] = currentEntry
		rmt.globalLock.Unlock()
	}

	return currentEntry.Get(current, rmt.calculator)
}

// rangeTreeEntry defines a helper struct for fast lookup of int range -> data.
type rangeTreeEntry struct {
	key          string
	internalTree *redblacktree.Tree
	entryLock    *sync.RWMutex
}

// Get returns the data found for the given range.
func (rte rangeTreeEntry) Get(current IntRange, calculator RangeMapTreeCalculator) interface{} {
	rte.entryLock.RLock()
	data, found := rte.internalTree.Get(current)
	rte.entryLock.RUnlock()

	if found {
		return data
	}

	resultRange, resultData := calculator(rte.key, current)

	rte.entryLock.Lock()
	rte.internalTree.Put(resultRange, resultData)
	rte.entryLock.Unlock()

	return resultData
}

func rangeComparator(a, b interface{}) int {
	i1 := a.(IntRange)
	i2 := b.(IntRange)

	if i1.StartPosition >= i2.StartPosition && i1.EndPosition <= i2.EndPosition {
		return 0
	}

	diff := i1.StartPosition - i2.StartPosition

	if diff < 0 {
		return -1
	}
	if diff > 0 {
		return 1
	}
	return 0
}
