// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate protoc --gofast_out=. membersig.proto

package proto

func (t *MemberSig) Name() string {
	return "MemberSig"
}

func (t *MemberSig) Value() string {
	bytes, err := t.Marshal()
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func (t *MemberSig) Build(value string) interface{} {
	uerr := t.Unmarshal([]byte(value))
	if uerr != nil {
		panic(uerr)
	}

	return t
}
