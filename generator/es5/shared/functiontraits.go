// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shared

// StateFunctionType defines an enumeration of types for a function.
type StateFunctionType int

const (
	StateFunctionNormalSync  StateFunctionType = iota
	StateFunctionNormalAsync StateFunctionType = iota
	StateFunctionSyncOrAsyncGenerator
)

// FunctionTraits creates and returns a StateFunctionTraits struct for the given traits of a function.
func FunctionTraits(isAsync bool, isGenerator bool, managesResources bool) StateFunctionTraits {
	return StateFunctionTraits{isAsync, isGenerator, managesResources}
}

// StateFunctionTraits is a simple struct for tracking the various traits of a function being generated.
type StateFunctionTraits struct {
	isAsync          bool
	isGenerator      bool
	managesResources bool
}

// IsAsynchronous returns true if the function is asynchronous, therefore returning a Promise.
func (sft StateFunctionTraits) IsAsynchronous() bool {
	return sft.isAsync
}

// IsGenerator returns true if a function is a generator, therefore returning a Stream<T>.
// (note that it can also be async, and therefore return a Promise<Stream<T>>)
func (sft StateFunctionTraits) IsGenerator() bool {
	return sft.isGenerator
}

// ManagesResources returns true if the function manages any resources under a `with` block.
func (sft StateFunctionTraits) ManagesResources() bool {
	return sft.managesResources
}

// Type returns an enumeration representing the aggregated traits for the function.
func (sft StateFunctionTraits) Type() StateFunctionType {
	switch {
	case sft.IsGenerator():
		return StateFunctionSyncOrAsyncGenerator

	case sft.IsAsynchronous():
		return StateFunctionNormalAsync

	default:
		return StateFunctionNormalSync
	}
}
