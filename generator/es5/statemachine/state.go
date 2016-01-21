// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"bytes"
)

// stateId represents a unique ID for a state.
type stateId int

// state represents a single state in the state machine.
type state struct {
	ID         stateId       // The ID of the current state.
	source     *bytes.Buffer // The source for this state.
	expression string        // The last expression value, if any.
	leafState  bool          // Whether this is a leaf state.
}

// Expression returns the expression value for this state, if any.
func (s *state) Expression() string {
	return s.expression
}

// Source returns the full source for this state.
func (s *state) Source() string {
	return s.source.String()
}

// HasSource returns true if the state has any source.
func (s *state) HasSource() bool {
	return s.source.Len() > 0
}

// IsLeafState returns true if the state is a leaf state in the state graph.
func (s *state) IsLeafState() bool {
	return s.leafState
}

// pushExpression adds the given expression to the current state.
func (s *state) pushExpression(expressionSource string) {
	s.expression = expressionSource
}

// pushSource adds source code to the current state.
func (s *state) pushSource(source string) {
	buf := s.source
	buf.WriteString(source)
	buf.WriteByte('\n')
}
