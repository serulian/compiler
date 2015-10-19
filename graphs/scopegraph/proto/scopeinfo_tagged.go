// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate protoc --gofast_out=. scopeinfo.proto

package proto

import (
	"github.com/serulian/compiler/graphs/typegraph"
)

func (t *ScopeInfo) Name() string {
	return "ScopeInfo"
}

func (t *ScopeInfo) Value() string {
	bytes, err := t.Marshal()
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func (t *ScopeInfo) Build(value string) interface{} {
	uerr := t.Unmarshal([]byte(value))
	if uerr != nil {
		panic(uerr)
	}

	return t
}

func (t *ScopeInfo) ReturnedTypeRef(tg *typegraph.TypeGraph) typegraph.TypeReference {
	if t.GetReturnedType() == "" {
		return tg.VoidTypeReference()
	}

	return tg.DeserializieTypeRef(t.GetReturnedType())
}
