// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

// GLOBAL_CONTEXT_ANNOTATIONS are the annotations that mark an interface as being a global context
// (e.g. Window) in WebIDL.
var GLOBAL_CONTEXT_ANNOTATIONS = []interface{}{"Global", "PrimaryGlobal"}

// CONSTRUCTOR_ANNOTATION is an annotation that describes support for a constructor on a WebIDL
// type. This translates to being able to do "new Type(...)" in ECMAScript.
const CONSTRUCTOR_ANNOTATION = "Constructor"

// NATIVE_OPERATOR_ANNOTATION is an annotation that marks an declaration as supporting the
// specified operator natively (i.e. not a custom defined operator).
const NATIVE_OPERATOR_ANNOTATION = "NativeOperator"

// SPECIALIZATION_NAMES maps WebIDL member specializations into Serulian typegraph names.
var SPECIALIZATION_NAMES = map[MemberSpecialization]string{
	GetterSpecialization: "index",
	SetterSpecialization: "setindex",
}

// NATIVE_TYPES maps from the predefined WebIDL types to the type actually supported
// in ES. We lose some information by doing so, but it allows for compatibility
// with existing WebIDL specifications. In the future, we might find a way to
// have these types be used in a more specific manner.
var NATIVE_TYPES = map[string]string{
	"boolean":             "Boolean",
	"byte":                "Number",
	"octet":               "Number",
	"short":               "Number",
	"unsigned short":      "Number",
	"long":                "Number",
	"unsigned long":       "Number",
	"long long":           "Number",
	"float":               "Number",
	"double":              "Number",
	"unrestricted float":  "Number",
	"unrestricted double": "Number",
}
