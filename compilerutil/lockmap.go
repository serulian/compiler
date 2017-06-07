// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

import "sync"

// LockMap defines a concurrent-safe map that returns a sync Mutex for each key. This is useful
// when multiple resources are being loaded concurrently and locking is needed, but only on a
// per-resource basis.
type LockMap struct {
	locks      map[string]*sync.Mutex
	globalLock *sync.Mutex
}

// CreateLockMap returns a new LockMap.
func CreateLockMap() LockMap {
	return LockMap{
		locks:      map[string]*sync.Mutex{},
		globalLock: &sync.Mutex{},
	}
}

// GetLock returns a lock for the given key.
func (lm LockMap) GetLock(key string) *sync.Mutex {
	lm.globalLock.Lock()
	defer lm.globalLock.Unlock()

	if lock, ok := lm.locks[key]; ok {
		return lock
	}

	lock := &sync.Mutex{}
	lm.locks[key] = lock
	return lock
}
