// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Stack copied from https://gist.github.com/bemasher/1777766

package parser

// errorProduction is a function that returns true if the given lexeme should *not* be
// consumed as part of error handling when consume fails. If *all* error productions on
// the stack return false, the token is consumed.
type errorProduction func(current commentedLexeme) bool

// errorProductionStack defines a stack for holding error production handlers.
type errorProductionStack struct {
	top  *errorProductionElement
	size int
}

type errorProductionElement struct {
	value errorProduction
	next  *errorProductionElement
}

// consumeToken executes the full set of error production handlers over the current token, return true if
// it should be consumed.
func (s *errorProductionStack) consumeToken(currentToken commentedLexeme) bool {
	if s.size == 0 {
		return true
	}

	var current = s.top
	for {
		if current == nil {
			return true
		}

		if current.value(currentToken) {
			return false
		}

		current = current.next
	}
}

// Push pushes an error production handler onto the stack.
func (s *errorProductionStack) push(value errorProduction) {
	s.top = &errorProductionElement{value, s.top}
	s.size++
}

// Pop removes the error production handler from the stack and returns it.
func (s *errorProductionStack) pop() (value errorProduction) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}
