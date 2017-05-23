// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"strings"

	"github.com/serulian/compiler/compilercommon"

	"fmt"

	cmap "github.com/streamrail/concurrent-map"
)

// SourceTracker is a helper struct for tracking the contents and versioning of all source files
// encountered by the package loader during its loading phase.
type SourceTracker struct {
	// sourceFiles is the map of tracked source files.
	sourceFiles map[compilercommon.InputSource]trackedSourceFile

	// pathLoader is the path loader used when loading the source files.
	pathLoader PathLoader
}

// ModifiedSourcePaths returns the source paths that have been modified since originally read,
// if any.
func (st SourceTracker) ModifiedSourcePaths() ([]compilercommon.InputSource, error) {
	var modifiedPaths = make([]compilercommon.InputSource, 0, len(st.sourceFiles))
	for _, file := range st.sourceFiles {
		current, err := st.pathLoader.GetRevisionID(string(file.sourcePath))
		if err != nil {
			return modifiedPaths, err
		}

		if current != file.revisionID {
			modifiedPaths = append(modifiedPaths, file.sourcePath)
		}
	}
	return modifiedPaths, nil
}

// LoadedContents returns the contents of the given source path when loaded by the packageloader, if any.
func (st SourceTracker) LoadedContents(path compilercommon.InputSource) ([]byte, bool) {
	tracked, exists := st.sourceFiles[path]
	if !exists {
		return []byte{}, false
	}

	return tracked.contents, true
}

func (st SourceTracker) RunePositionToLineAndCol(runePosition uint64, path compilercommon.InputSource) (uint64, uint64, error) {
	tsf, exists := st.sourceFiles[path]
	if !exists {
		return 0, 0, fmt.Errorf("Could not find path %s", path)
	}

	if tsf.positionMapper == nil {
		tsf.positionMapper = createSourcePositionMapper(tsf.contents)
	}

	return tsf.positionMapper.RunePositionToLineAndCol(runePosition)
}

func (st SourceTracker) LineAndColToRunePosition(lineNumber uint64, colPosition uint64, path compilercommon.InputSource) (uint64, error) {
	tsf, exists := st.sourceFiles[path]
	if !exists {
		return 0, fmt.Errorf("Could not find path %s", path)
	}

	if tsf.positionMapper == nil {
		tsf.positionMapper = createSourcePositionMapper(tsf.contents)
	}

	return tsf.positionMapper.LineAndColToRunePosition(lineNumber, colPosition)
}

func (st SourceTracker) TextForLine(lineNumber uint64, path compilercommon.InputSource) (string, error) {
	tsf, exists := st.sourceFiles[path]
	if !exists {
		return "", fmt.Errorf("Could not find path %s", path)
	}

	lines := strings.Split(string(tsf.contents), "\n")
	if int(lineNumber) >= len(lines) {
		return "", fmt.Errorf("Line number %v not found in path %v", lineNumber, path)
	}

	return lines[int(lineNumber)], nil
}

// trackedSourceFile defines a struct for tracking the versioning and contents of a source file
// encountered during package loading.
type trackedSourceFile struct {
	// sourcePath is the path of the source file encountered.
	sourcePath compilercommon.InputSource

	// sourceKind is the kind of the source file.
	sourceKind string

	// contents are the contents of the source file when it was loaded.
	contents []byte

	// revisionID is the revision ID of the source file when it was loaded, as returned by
	// GetRevisionID on the PathLoader. Typically, this is an mtime or version number. If the
	// revision ID has changed, the contents of the file are assumed to have been altered since
	// the file was read.
	revisionID int64

	// positionMapper is the source position mapper for converting rune positions <-> line+col
	// positions. Only instantiated if necessary.
	positionMapper *sourcePositionMapper
}

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
	m.sourceFiles.Set(string(path), trackedSourceFile{
		sourcePath: path,
		sourceKind: sourceKind,
		contents:   contents,
		revisionID: revisionID,
	})
}

func (m *mutableSourceTracker) RunePositionToLineAndCol(runePosition uint64, path compilercommon.InputSource) (uint64, uint64, error) {
	tracked, exists := m.sourceFiles.Get(string(path))
	if !exists {
		return 0, 0, fmt.Errorf("Could not find path %s", path)
	}

	tsf := tracked.(trackedSourceFile)
	if tsf.positionMapper == nil {
		tsf.positionMapper = createSourcePositionMapper(tsf.contents)
	}

	return tsf.positionMapper.RunePositionToLineAndCol(runePosition)
}

func (m *mutableSourceTracker) LineAndColToRunePosition(lineNumber uint64, colPosition uint64, path compilercommon.InputSource) (uint64, error) {
	tracked, exists := m.sourceFiles.Get(string(path))
	if !exists {
		return 0, fmt.Errorf("Could not find path %s", path)
	}

	tsf := tracked.(trackedSourceFile)
	if tsf.positionMapper == nil {
		tsf.positionMapper = createSourcePositionMapper(tsf.contents)
	}

	return tsf.positionMapper.LineAndColToRunePosition(lineNumber, colPosition)
}

func (m *mutableSourceTracker) TextForLine(lineNumber uint64, path compilercommon.InputSource) (string, error) {
	tracked, exists := m.sourceFiles.Get(string(path))
	if !exists {
		return "", fmt.Errorf("Could not find path %s", path)
	}

	tsf := tracked.(trackedSourceFile)

	lines := strings.Split(string(tsf.contents), "\n")
	if int(lineNumber) >= len(lines) {
		return "", fmt.Errorf("Line number %v not found in path %v", lineNumber, path)
	}

	return lines[int(lineNumber)], nil
}

// Freeze freezes the mutable source tracker into an immutable SourceTracker.
func (m *mutableSourceTracker) Freeze() SourceTracker {
	sourceFileMap := map[compilercommon.InputSource]trackedSourceFile{}
	for entry := range m.sourceFiles.Iter() {
		sourceFileMap[compilercommon.InputSource(entry.Key)] = entry.Val.(trackedSourceFile)
	}
	return SourceTracker{sourceFileMap, m.pathLoader}
}
