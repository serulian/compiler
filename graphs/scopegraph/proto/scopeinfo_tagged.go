// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate protoc --gofast_out=. scopeinfo.proto

package proto

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
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

func (t *ScopeInfo) CalledOperator(tg *typegraph.TypeGraph) (typegraph.TGMember, bool) {
	if t.CalledOpReference == nil {
		return typegraph.TGMember{}, false
	}

	nodeId := compilergraph.GraphNodeId(t.CalledOpReference.GetReferencedNode())
	return tg.GetTypeOrMember(nodeId).(typegraph.TGMember), true
}

func (t *ScopeInfo) NamedReferenceNode(srg *srg.SRG, tg *typegraph.TypeGraph) (compilergraph.GraphNode, bool) {
	if t.NamedReference == nil {
		return compilergraph.GraphNode{}, false
	}

	nodeId := compilergraph.GraphNodeId(t.NamedReference.GetReferencedNode())
	if t.NamedReference.GetIsSRGNode() {
		return srg.GetNode(nodeId), true
	} else {
		return tg.GetNode(nodeId), true
	}
}

func (t *ScopeInfo) Build(value string) interface{} {
	uerr := t.Unmarshal([]byte(value))
	if uerr != nil {
		panic(uerr)
	}

	return t
}

func (t *ScopeInfo) StaticTypeRef(tg *typegraph.TypeGraph) typegraph.TypeReference {
	if t.GetStaticType() == "" {
		return tg.VoidTypeReference()
	}

	return tg.DeserializieTypeRef(t.GetStaticType())
}

func (t *ScopeInfo) AssignableTypeRef(tg *typegraph.TypeGraph) typegraph.TypeReference {
	if t.GetAssignableType() == "" {
		return tg.VoidTypeReference()
	}

	return tg.DeserializieTypeRef(t.GetAssignableType())
}

func (t *ScopeInfo) ResolvedTypeRef(tg *typegraph.TypeGraph) typegraph.TypeReference {
	if t.GetResolvedType() == "" {
		return tg.VoidTypeReference()
	}

	return tg.DeserializieTypeRef(t.GetResolvedType())
}

func (t *ScopeInfo) ReturnedTypeRef(tg *typegraph.TypeGraph) typegraph.TypeReference {
	if t.GetReturnedType() == "" {
		return tg.VoidTypeReference()
	}

	return tg.DeserializieTypeRef(t.GetReturnedType())
}
