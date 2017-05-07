// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"os"
	"sync"
)

// From: http://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-denoted-by-a-path-exists-in-golang
// exists returns whether the given file or directory exists or not.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// lockMap defines a concurrent-safe map that returns a sync Mutex for each key. This is useful
// when multiple resources are being loaded concurrently and locking is needed, but only on a
// per-resource basis.
type lockMap struct {
	locks      map[string]*sync.Mutex
	globalLock *sync.Mutex
}

// createLockMap returns a new lockMap.
func createLockMap() lockMap {
	return lockMap{
		locks:      map[string]*sync.Mutex{},
		globalLock: &sync.Mutex{},
	}
}

// getLock returns a lock for the given key.
func (lm lockMap) getLock(key string) *sync.Mutex {
	lm.globalLock.Lock()
	defer lm.globalLock.Unlock()

	if lock, ok := lm.locks[key]; ok {
		return lock
	}

	lock := &sync.Mutex{}
	lm.locks[key] = lock
	return lock
}
