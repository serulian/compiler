// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codedom

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph"
)

// RuntimeFunction defines a function defined by the runtime.
type RuntimeFunction string

const (
	CastFunction               RuntimeFunction = "$t.cast"
	IsTypeFunction             RuntimeFunction = "$t.istype"
	DynamicAccessFunction      RuntimeFunction = "$t.dynamicaccess"
	AssertNotNullFunction      RuntimeFunction = "$t.assertnotnull"
	StreamMemberAccessFunction RuntimeFunction = "$t.streamaccess"
	BoxFunction                RuntimeFunction = "$t.box"
	FastBoxFunction            RuntimeFunction = "$t.fastbox"
	UnboxFunction              RuntimeFunction = "$t.unbox"
	NullableInvokeFunction     RuntimeFunction = "$t.nullableinvoke"

	AsyncNullableComparisonFunction RuntimeFunction = "$t.asyncnullcompare"
	SyncNullableComparisonFunction  RuntimeFunction = "$t.syncnullcompare"

	NewPromiseFunction          RuntimeFunction = "$promise.new"
	ResolvePromiseFunction      RuntimeFunction = "$promise.resolve"
	TranslatePromiseFunction    RuntimeFunction = "$promise.translate"
	ShortCircuitPromiseFunction RuntimeFunction = "$promise.shortcircuit"
	MaybePromiseFunction        RuntimeFunction = "$promise.maybe"

	StatePushResourceFunction RuntimeFunction = "$resources.pushr"
	StatePopResourceFunction  RuntimeFunction = "$resources.popr"

	EmptyGeneratorDirect RuntimeFunction = "$generator.directempty"

	BoxedDataProperty string = "$wrapped"
)

// RuntimeFunctionCallNode represents a call to an internal runtime function defined
// for special handling of code.
type RuntimeFunctionCallNode struct {
	expressionBase
	Function  RuntimeFunction // The runtime function being called.
	Arguments []Expression    // The arguments for the call.
}

func RuntimeFunctionCall(function RuntimeFunction, arguments []Expression, basis compilergraph.GraphNode) Expression {
	return &RuntimeFunctionCallNode{
		expressionBase{domBase{basis}},
		function,
		arguments,
	}
}

func (e *RuntimeFunctionCallNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return isAsynchronous(scopegraph, e.Arguments)
}
