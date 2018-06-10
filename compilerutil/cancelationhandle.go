// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

// CancelFunction is a function that can be invoked to cancel the operation.
type CancelFunction func()

// CancelationHandle defines a handle for the cancelation of operations.
type CancelationHandle interface {
	// Cancel marks the operation as canceled.
	Cancel()

	// WasCanceled returns whether the operation was canceled.
	WasCanceled() bool
}

// NewCancelationHandle returns a new handle for canceling an operation.
func NewCancelationHandle() CancelationHandle {
	return &cancelationHandle{}
}

// NoopCancelationHandle returns a cancelation handle that cannot be canceled.
func NoopCancelationHandle() CancelationHandle {
	return noopCancelationHandle{}
}

// GetCancelationHandle returns either the existing cancelation handle (if not nil)
// or a new no-op handle.
func GetCancelationHandle(existing CancelationHandle) CancelationHandle {
	if existing == nil {
		return noopCancelationHandle{}
	}

	return existing
}

type cancelationHandle struct {
	canceled bool
}

func (c *cancelationHandle) WasCanceled() bool {
	return c.canceled
}

func (c *cancelationHandle) Cancel() {
	c.canceled = true
}

type noopCancelationHandle struct{}

func (noopCancelationHandle) WasCanceled() bool {
	return false
}

func (noopCancelationHandle) Cancel() {
	panic("Should never be called")
}
