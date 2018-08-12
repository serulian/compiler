// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"strings"

	"github.com/serulian/compiler/compilercommon"

	"fmt"

	diffmatchpatch "github.com/sergi/go-diff/diffmatchpatch"
)

// PositionType defines the type of positions given to GetPositionOffsetFromTracked.
type PositionType int

const (
	// TrackedFilePosition indicates a position in the tracked file contents.
	TrackedFilePosition PositionType = 1

	// CurrentFilePosition indicates a position in the current file contents.
	CurrentFilePosition = -1
)

// SourceTracker is a helper struct for tracking the contents and versioning of all source files
// encountered by the package loader.
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

// RunePositionToLineAndCol returns the line number and column position of the rune position in the source of the given path.
func (st SourceTracker) RunePositionToLineAndCol(runePosition int, path compilercommon.InputSource, sourceOption compilercommon.SourceMappingOption) (int, int, error) {
	positionMapper, err := st.getPositionMapper(path, sourceOption)
	if err != nil {
		return 0, 0, err
	}
	return positionMapper.RunePositionToLineAndCol(runePosition)
}

// LineAndColToRunePosition returns the rune position of the line number and column position in the source of the given path.
func (st SourceTracker) LineAndColToRunePosition(lineNumber int, colPosition int, path compilercommon.InputSource, sourceOption compilercommon.SourceMappingOption) (int, error) {
	positionMapper, err := st.getPositionMapper(path, sourceOption)
	if err != nil {
		return 0, err
	}

	return positionMapper.LineAndColToRunePosition(lineNumber, colPosition)
}

// TextForLine returns all the text found on the given (0-indexed) line in the source of the given path.
func (st SourceTracker) TextForLine(lineNumber int, path compilercommon.InputSource, sourceOption compilercommon.SourceMappingOption) (string, error) {
	contents, err := st.getContents(path, sourceOption)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(contents), "\n")
	if int(lineNumber) >= len(lines) {
		return "", fmt.Errorf("Line number %v not found in path %v with option %v", lineNumber, path, sourceOption)
	}

	return lines[int(lineNumber)], nil
}

// getPositionMapper returns a position mapper for the path, caching as necessary.
func (st SourceTracker) getPositionMapper(path compilercommon.InputSource, sourceOption compilercommon.SourceMappingOption) (compilercommon.SourcePositionMapper, error) {
	tsf, exists := st.sourceFiles[path]
	if !exists {
		return compilercommon.EmptySourcePositionMapper(), fmt.Errorf("Could not find path %s", path)
	}

	switch sourceOption {
	case compilercommon.SourceMapCurrent:
		currentVersion, err := st.pathLoader.GetRevisionID(string(path))
		if err != nil {
			return compilercommon.EmptySourcePositionMapper(), err
		}

		latest := tsf.latestPositionMapper
		if latest == nil || latest.revisionID != currentVersion {
			currentContents, err := st.getContents(path, sourceOption)
			if err != nil {
				return compilercommon.EmptySourcePositionMapper(), err
			}

			latest = &positionMapperAndRevisionID{currentVersion, compilercommon.CreateSourcePositionMapper(currentContents)}
			tsf.latestPositionMapper = latest
		}

		return latest.positionMapper, nil

	case compilercommon.SourceMapTracked:
		tracked := tsf.trackedPositionMapper
		if tracked == nil {
			tracked = &positionMapperAndRevisionID{tsf.revisionID, compilercommon.CreateSourcePositionMapper(tsf.contents)}
			tsf.trackedPositionMapper = tracked
		}

		return tracked.positionMapper, nil

	default:
		panic("Unknown source option")
	}
}

// getContents returns the contents of the path.
func (st SourceTracker) getContents(path compilercommon.InputSource, sourceOption compilercommon.SourceMappingOption) ([]byte, error) {
	tsf, exists := st.sourceFiles[path]
	if !exists {
		return []byte{}, fmt.Errorf("Could not find path %s", path)
	}

	currentVersion, err := st.pathLoader.GetRevisionID(string(path))
	if err != nil {
		return []byte{}, err
	}

	if sourceOption != compilercommon.SourceMapCurrent || currentVersion == tsf.revisionID {
		return tsf.contents, nil
	}

	return st.pathLoader.LoadSourceFile(string(path))
}
