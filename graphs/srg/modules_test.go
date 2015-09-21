// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"testing"

	"github.com/serulian/compiler/compilergraph"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func TestBasicModules(t *testing.T) {
	graph, err := compilergraph.NewGraph("tests/basic/basic.seru")
	if err != nil {
		t.Errorf("%v", err)
	}

	testSRG := NewSRG(graph)
	result := testSRG.LoadAndParse()
	if !result.Status {
		t.Errorf("Expected successful parse")
	}

	// Ensure that both modules were loaded.
	modules := testSRG.GetModules()

	if len(modules) != 2 {
		t.Errorf("Expected 2 modules found, found: %v", modules)
	}

	var modulePaths []string
	for _, module := range modules {
		modulePaths = append(modulePaths, string(module.InputSource))
		node, found := testSRG.FindModuleBySource(module.InputSource)

		assert.True(t, found, "Could not find module %s", module.InputSource)
		assert.Equal(t, node.fileNode.NodeId, module.fileNode.NodeId, "Node ID mismatch on modules")
	}

	assert.Contains(t, modulePaths, "tests/basic/basic.seru", "Missing basic.seru module")
	assert.Contains(t, modulePaths, "tests/basic/anotherfile.seru", "Missing anotherfile.seru module")
}
