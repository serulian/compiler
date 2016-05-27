// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"github.com/serulian/compiler/generator/escommon/esbuilder"
)

// stateId represents a unique ID for a state.
type stateId int

// state represents a single state in the state machine.
type state struct {
	ID        stateId                   // The ID of the current state.
	pieces    []esbuilder.SourceBuilder // The pieces of the generated source for the state.
	leafState bool                      // Whether this is a leaf state.
}

// Source returns the full source for this state.
func (s *state) SourceTemplate() esbuilder.SourceBuilder {
	templateStr := `{{ range $piece := .Pieces }}
{{ emit $piece }}
{{ end }}`

	return esbuilder.Template("state", templateStr, s)
}

// Pieces returns the builder pieces add to this state.
func (s *state) Pieces() []esbuilder.SourceBuilder {
	return s.pieces
}

// HasSource returns true if the state has any source.
func (s *state) HasSource() bool {
	return len(s.pieces) > 0
}

// IsLeafState returns true if the state is a leaf state in the state graph.
func (s *state) IsLeafState() bool {
	return s.leafState
}

// pushSnippet adds source code to the current state.
func (s *state) pushSnippet(source string) {
	s.pushBuilt(esbuilder.Snippet(source))
}

// pushBuilt adds builder to the current state.
func (s *state) pushBuilt(builder esbuilder.SourceBuilder) {
	s.pieces = append(s.pieces, builder)
}
