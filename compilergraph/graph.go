// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"github.com/cayleygraph/cayley"

	"github.com/serulian/compiler/compilerutil"
)

// serulianGraph represents a full Serulian graph, including its AST (SRG), type system (TG)
// and other various graphs used by the compiler.
type serulianGraph struct {
	rootSourceFilePath string         // The root source file path.
	cayleyStore        *cayley.Handle // Handle to the cayley store.
}

// NewGraphLayer returns a new graph layer of the given kind.
func (sg *serulianGraph) NewGraphLayer(uniqueID string, nodeKindEnum TaggedValue) GraphLayer {
	return &graphLayer{
		id:                compilerutil.NewUniqueId(),
		prefix:            uniqueID,
		cayleyStore:       sg.cayleyStore,
		nodeKindPredicate: "node-kind",
		nodeKindEnum:      nodeKindEnum,
	}
}

func (sg *serulianGraph) RootSourceFilePath() string {
	return sg.rootSourceFilePath
}
