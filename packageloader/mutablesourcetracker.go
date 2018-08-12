// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"strings"

	"github.com/serulian/compiler/compilercommon"
	cmap "github.com/streamrail/concurrent-map"

	"fmt"
)

// mutableSourceTracker is a source tracker that is built during the package loading phase.
type mutableSourceTracker struct {
	// sourceFiles is a map from source file path to a trackedSourceFile struct.
	sourceFiles cmap.ConcurrentMap

	// pathLoader is the path loader used when loading the source files.
	pathLoader PathLoader
}

func newMutableSourceTracker(pathLoader PathLoader) *mutableSourceTracker {
	return &mutableSourceTracker{
		sourceFiles: cmap.New(),
		pathLoader:  pathLoader,
	}
}

// AddSourceFile adds a source file to the mutable source tracker.
func (m *mutableSourceTracker) AddSourceFile(path compilercommon.InputSource, sourceKind string, contents []byte, revisionID int64) {
	m.sourceFiles.Set(string(path), &trackedSourceFile{
		sourcePath: path,
		sourceKind: sourceKind,
		contents:   contents,
		revisionID: revisionID,
	})
}

func (m *mutableSourceTracker) RunePositionToLineAndCol(runePosition int, path compilercommon.InputSource, sourceOption compilercommon.SourceMappingOption) (int, int, error) {
	if sourceOption != compilercommon.SourceMapTracked {
		panic("Cannot load non-tracked source in mutable source tracker")
	}

	tracked, exists := m.sourceFiles.Get(string(path))
	if !exists {
		return 0, 0, fmt.Errorf("Could not find path %s", path)
	}

	tsf := tracked.(*trackedSourceFile)
	return m.getPositionMapper(tsf).RunePositionToLineAndCol(runePosition)
}

func (m *mutableSourceTracker) LineAndColToRunePosition(lineNumber int, colPosition int, path compilercommon.InputSource, sourceOption compilercommon.SourceMappingOption) (int, error) {
	if sourceOption != compilercommon.SourceMapTracked {
		panic("Cannot load non-tracked source in mutable source tracker")
	}

	tracked, exists := m.sourceFiles.Get(string(path))
	if !exists {
		return 0, fmt.Errorf("Could not find path %s", path)
	}

	tsf := tracked.(*trackedSourceFile)
	return m.getPositionMapper(tsf).LineAndColToRunePosition(lineNumber, colPosition)
}

func (m *mutableSourceTracker) TextForLine(lineNumber int, path compilercommon.InputSource, sourceOption compilercommon.SourceMappingOption) (string, error) {
	if sourceOption != compilercommon.SourceMapTracked {
		panic("Cannot load non-tracked source in mutable source tracker")
	}

	tracked, exists := m.sourceFiles.Get(string(path))
	if !exists {
		return "", fmt.Errorf("Could not find path %s", path)
	}

	tsf := tracked.(*trackedSourceFile)

	lines := strings.Split(string(tsf.contents), "\n")
	if int(lineNumber) >= len(lines) {
		return "", fmt.Errorf("Line number %v not found in path %v", lineNumber, path)
	}

	return lines[int(lineNumber)], nil
}

// Freeze freezes the mutable source tracker into an immutable SourceTracker.
func (m *mutableSourceTracker) Freeze() SourceTracker {
	sourceFileMap := map[compilercommon.InputSource]*trackedSourceFile{}
	for entry := range m.sourceFiles.Iter() {
		sourceFileMap[compilercommon.InputSource(entry.Key)] = entry.Val.(*trackedSourceFile)
	}
	return SourceTracker{sourceFileMap, m.pathLoader}
}

func (m *mutableSourceTracker) getPositionMapper(tsf *trackedSourceFile) compilercommon.SourcePositionMapper {
	if tsf.trackedPositionMapper == nil {
		tsf.trackedPositionMapper = &positionMapperAndRevisionID{tsf.revisionID, compilercommon.CreateSourcePositionMapper(tsf.contents)}
	}
	return tsf.trackedPositionMapper.positionMapper
}
