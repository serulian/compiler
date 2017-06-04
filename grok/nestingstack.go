// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Stack copied from https://gist.github.com/bemasher/1777766

package grok

type nestingStack struct {
	top  *nestingElement
	size int
}

type nestingElement struct {
	openRune           rune
	startIndex         int
	commaCounter       int
	lastSeparatorIndex int
	next               *nestingElement
}

// expressionStartIndex returns the start index for the nested expression, based on the
// start index of the nesting element and last encountered separator index.
func (ne *nestingElement) expressionStartIndex() int {
	startIndex := ne.startIndex + 1
	if ne.lastSeparatorIndex+1 > startIndex {
		return ne.lastSeparatorIndex + 1
	}
	return startIndex
}

func (s *nestingStack) topValue() (rune, int, bool) {
	if s.size == 0 {
		return ' ', -1, false
	}

	return s.top.openRune, s.top.startIndex, true
}

// Push pushes an opening rune onto the stack.
func (s *nestingStack) push(openRune rune, startIndex int) {
	s.top = &nestingElement{openRune, startIndex, 0, -1, s.top}
	s.size++
}

// Pop removes the nesting rune from the stack and returns it.
func (s *nestingStack) pop() (openRune rune, startIndex int) {
	if s.size > 0 {
		openRune, startIndex, s.top = s.top.openRune, s.top.startIndex, s.top.next
		s.size--
		return
	}
	return ' ', -1
}
