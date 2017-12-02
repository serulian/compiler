// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/sourceshape"

	"github.com/stretchr/testify/assert"
)

func TestFindNodeForPosition(t *testing.T) {
	testSRG := getSRG(t, "tests/position/position.seru")
	cit := testSRG.AllComments()
	for cit.Next() {
		comment := SRGComment{cit.Node(), testSRG}
		parent := comment.ParentNode()

		source := parent.Get(sourceshape.NodePredicateSource)
		startRune := parent.GetValue(sourceshape.NodePredicateStartRune).Int()

		sourcePosition := compilercommon.InputSource(source).PositionForRunePosition(startRune, nil)
		node, found := testSRG.FindNodeForPosition(sourcePosition)
		if !assert.True(t, found, "Missing node with comment %s", comment.Contents()) {
			continue
		}

		if !assert.Equal(t, node.GetNodeId(), parent.GetNodeId(), "Mismatch of node with comment %s", comment.Contents()) {
			continue
		}
	}
}
