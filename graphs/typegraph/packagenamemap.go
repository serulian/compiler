// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"

	cmap "github.com/streamrail/concurrent-map"
)

type packagenamemap struct {
	internalmap cmap.ConcurrentMap
}

func newPackageNameMap() packagenamemap {
	return packagenamemap{
		internalmap: cmap.New(),
	}
}

// CheckAndTrack adds the given type or member to the package map, if it is the first
// instance seen. If it is not, then false and the existing type or member, under the
// package, with the same name is returned.
func (pnm packagenamemap) CheckAndTrack(module TGModule, typeOrMember TGTypeOrMember) (TGTypeOrMember, bool) {
	key := fmt.Sprintf("%v::%v::%v", module.SourceGraphId(), module.PackagePath(), typeOrMember.Name())
	if !pnm.internalmap.SetIfAbsent(key, typeOrMember) {
		existing, _ := pnm.internalmap.Get(key)
		return existing.(TGTypeOrMember), false
	}

	return nil, true
}
