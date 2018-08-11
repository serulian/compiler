// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"strings"

	"github.com/serulian/compiler/compilercommon"

	"fmt"

	diffmatchpatch "github.com/sergi/go-diff/diffmatchpatch"
	cmap "github.com/streamrail/concurrent-map"
)

// SourceTracker is a helper struct for tracking the contents and versioning of all source files
// encountered by the package loader during its loading phase.
type SourceTracker struct {
	// sourceFiles is the map of tracked source files.
	sourceFiles map[compilercommon.InputSource]*trackedSourceFile

	// pathLoader is the path loader used when loading the source files.
	pathLoader PathLoader
}

// HasModifiedSourcePaths returns whether any of the source paths have been modified since originally read.
func (st SourceTracker) HasModifiedSourcePaths() (bool, error) {
	for _, file := range st.sourceFiles {
		current, err := st.pathLoader.GetRevisionID(string(file.sourcePath))
		if err != nil {
			return false, err
		}

		if current != file.revisionID {
			return true, nil
		}
	}
	return false, nil
}

// PositionType defines the type of positions given to GetPositionOffsetFromTracked.
type PositionType int

const (
	// TrackedFilePosition indicates a position in the tracked file contents.
	TrackedFilePosition PositionType = 1

	// CurrentFilePosition indicates a position in the current file contents.
	CurrentFilePosition = -1
)

// GetPositionOffset returns the given position, offset by any changes that occured since the source file
// referenced in the position has been tracked.
func (st SourceTracker) GetPositionOffset(position compilercommon.SourcePosition, positionType PositionType) (compilercommon.SourcePosition, error) {
	diffs, err := st.DiffFromTracked(position.Source())
	if err != nil {
		return position, err
	}

	if len(diffs) == 0 {
		return position, nil
	}

	// Otherwise, adjust the position using the diffs.
	runePosition, err := position.RunePosition()
	if err != nil {
		return position, err
	}

	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffEqual:
			break

		case diffmatchpatch.DiffInsert:
			runePosition += (len(diff.Text) * int(positionType))

		case diffmatchpatch.DiffDelete:
			runePosition -= (len(diff.Text) * int(positionType))

		default:
			panic("Unknown kind of diff")
		}
	}

	if runePosition < 0 {
		runePosition = 0
	}

	return position.Source().PositionForRunePosition(runePosition, st), nil
}

// DiffFromTracked returns the diff of the current version of the path (as loaded via the path
// loader) from the contents stored in the tracker.
func (st SourceTracker) DiffFromTracked(path compilercommon.InputSource) ([]diffmatchpatch.Diff, error) {
	tracked, exists := st.sourceFiles[path]
	if !exists {
		return []diffmatchpatch.Diff{}, fmt.Errorf("Path specified is not being tracked")
	}

	currentVersion, err := st.pathLoader.GetRevisionID(string(path))
	if err != nil {
		return []diffmatchpatch.Diff{}, err
	}

	if currentVersion == tracked.revisionID {
		// No changes.
		return []diffmatchpatch.Diff{}, err
	}

	// Otherwise, load the current contents and diff.
	currentContents, err := st.pathLoader.LoadSourceFile(string(path))
	if err != nil {
		return []diffmatchpatch.Diff{}, err
	}

	dmp := diffmatchpatch.New()
	return dmp.DiffCleanupSemantic(dmp.DiffMain(string(tracked.contents), string(currentContents), false)), nil
}

// LoadedContents returns the contents of the given source path when loaded by the packageloader, if any.
func (st SourceTracker) LoadedContents(path compilercommon.InputSource) ([]byte, bool) {
	tracked, exists := st.sourceFiles[path]
	if !exists {
		return []byte{}, false
	}

	return tracked.contents, true
}

func (st SourceTracker) RunePositionToLineAndCol(runePosition int, path compilercommon.InputSource, sourceOption compilercommon.SourceMappingOption) (int, int, error) {
	tsf, exists := st.sourceFiles[path]
	if !exists {
		return 0, 0, fmt.Errorf("Could not find path %s", path)
	}

	if sourceOption == compilercommon.SourceMapCurrent {
		currentContents, err := st.pathLoader.LoadSourceFile(string(path))
		if err != nil {
			return 0, 0, err
		}

		return compilercommon.CreateSourcePositionMapper(currentContents).RunePositionToLineAndCol(runePosition)
	}

	if tsf.positionMapper == nil {
		tsf.positionMapper = compilercommon.CreateSourcePositionMapper(tsf.contents)
	}

	return tsf.positionMapper.RunePositionToLineAndCol(runePosition)
}

func (st SourceTracker) LineAndColToRunePosition(lineNumber int, colPosition int, path compilercommon.InputSource, sourceOption compilercommon.SourceMappingOption) (int, error) {
	tsf, exists := st.sourceFiles[path]
	if !exists {
		return 0, fmt.Errorf("Could not find path %s", path)
	}

	if sourceOption == compilercommon.SourceMapCurrent {
		currentContents, err := st.pathLoader.LoadSourceFile(string(path))
		if err != nil {
			return 0, err
		}

		return compilercommon.CreateSourcePositionMapper(currentContents).LineAndColToRunePosition(lineNumber, colPosition)
	}

	if tsf.positionMapper == nil {
		tsf.positionMapper = compilercommon.CreateSourcePositionMapper(tsf.contents)
	}

	return tsf.positionMapper.LineAndColToRunePosition(lineNumber, colPosition)
}

func (st SourceTracker) TextForLine(lineNumber int, path compilercommon.InputSource, sourceOption compilercommon.SourceMappingOption) (string, error) {
	tsf, exists := st.sourceFiles[path]
	if !exists {
		return "", fmt.Errorf("Could not find path %s", path)
	}

	contents := tsf.contents
	if sourceOption == compilercommon.SourceMapCurrent {
		currentContents, err := st.pathLoader.LoadSourceFile(string(path))
		if err != nil {
			return "", err
		}

		contents = currentContents
	}

	lines := strings.Split(string(contents), "\n")
	if int(lineNumber) >= len(lines) {
		return "", fmt.Errorf("Line number %v not found in path %v with option %v", lineNumber, path, sourceOption)
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
	positionMapper *compilercommon.SourcePositionMapper
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
	if tsf.positionMapper == nil {
		tsf.positionMapper = compilercommon.CreateSourcePositionMapper(tsf.contents)
	}

	return tsf.positionMapper.RunePositionToLineAndCol(runePosition)
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
	if tsf.positionMapper == nil {
		tsf.positionMapper = compilercommon.CreateSourcePositionMapper(tsf.contents)
	}

	return tsf.positionMapper.LineAndColToRunePosition(lineNumber, colPosition)
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
