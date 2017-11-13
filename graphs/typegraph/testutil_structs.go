// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import "fmt"

// TestModule defines a new mdoule to be created in the test type graph.
type TestModule struct {
	ModuleName string
	Types      []TestType
	Members    []TestMember
}

func (tm TestModule) key() string {
	return fmt.Sprintf("module-%v", tm.ModuleName)
}

// TestType defines a new type to be created in the test type graph.
type TestType struct {
	Kind       string
	Name       string
	ParentType string
	Generics   []TestGeneric
	Members    []TestMember
}

func (tt TestType) key() string {
	return fmt.Sprintf("type-%v", tt.Name)
}

// TestGeneric defines a new generic to be created on a TestType or TestMember.
type TestGeneric struct {
	Name       string
	Constraint string
}

func (tg TestGeneric) key() string {
	return fmt.Sprintf("generic-%v", tg.Name)
}

// TestMember defines a new member to be created on a TestType or under a test module.
type TestMember struct {
	Kind       MemberSignatureKind
	Name       string
	ReturnType string
	Generics   []TestGeneric
	Parameters []TestParam
}

func (tm TestMember) key() string {
	return fmt.Sprintf("member-%v", tm.Name)
}

// TestParam defines a new parameter to be created on a TestMember.
type TestParam struct {
	Name      string
	ParamType string
}

func (tp TestParam) key() string {
	return fmt.Sprintf("param-%v", tp.Name)
}
