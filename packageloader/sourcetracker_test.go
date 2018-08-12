// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"fmt"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilercommon"

	"github.com/stretchr/testify/assert"
)

type anotherTestFile struct {
	contents   []byte
	revisionID int64
}

type anotherTestPathLoader struct {
	files map[string]anotherTestFile
}

func (tpl *anotherTestPathLoader) setFile(path string, contents string) {
	var revisionID int64
	file, exists := tpl.files[path]
	if exists {
		revisionID = file.revisionID + 1
	}

	tpl.files[path] = anotherTestFile{
		contents:   []byte(contents),
		revisionID: revisionID,
	}
}

func (tpl *anotherTestPathLoader) Exists(path string) (bool, error) {
	_, exists := tpl.files[path]
	return exists, nil
}

func (tpl *anotherTestPathLoader) LoadSourceFile(path string) ([]byte, error) {
	file, exists := tpl.files[path]
	if exists {
		return file.contents, nil
	}

	return []byte{}, fmt.Errorf("Could not find file: %s", path)
}

func (tpl *anotherTestPathLoader) IsSourceFile(path string) bool {
	_, exists := tpl.files[path]
	return exists
}

func (tpl *anotherTestPathLoader) LoadDirectory(path string) ([]DirectoryEntry, error) {
	return []DirectoryEntry{}, fmt.Errorf("Invalid path: %s", path)
}

func (tpl *anotherTestPathLoader) VCSPackageDirectory(entrypoint Entrypoint) string {
	return ""
}

func (tpl *anotherTestPathLoader) GetRevisionID(path string) (int64, error) {
	file, exists := tpl.files[path]
	if exists {
		return file.revisionID, nil
	}

	return 0, fmt.Errorf("Could not find file: %s", path)
}

func assertPositionMapping(t *testing.T, tracker SourceTracker, sourceInput compilercommon.InputSource, sourceText string, sourceOption compilercommon.SourceMappingOption) {
	lines := strings.Split(sourceText, "\n")

	var counter = 0
	for lineIndex, lineText := range lines {
		for charIndex := range lineText {
			l, c, err := tracker.RunePositionToLineAndCol(counter+charIndex, sourceInput, sourceOption)
			assert.Nil(t, err)
			assert.Equal(t, lineIndex, l)
			assert.Equal(t, charIndex, c)

			r, err := tracker.LineAndColToRunePosition(lineIndex, charIndex, sourceInput, sourceOption)
			assert.Nil(t, err)
			assert.Equal(t, r, counter+charIndex)
		}
		counter += len(lineText) + 1
	}
}

func TestSourceTracker(t *testing.T) {
	pathLoader := &anotherTestPathLoader{
		files: map[string]anotherTestFile{},
	}

	fooContents := "some\nfoo\nfile"
	pathLoader.setFile("foo.txt", fooContents)

	fooSource := compilercommon.InputSource("foo.txt")

	mutableTracker := newMutableSourceTracker(pathLoader)
	mutableTracker.AddSourceFile(fooSource, "", []byte(fooContents), 0)

	// Check tracked contents.
	tracker := mutableTracker.Freeze()
	trackedFooContents, _ := tracker.LoadedContents(fooSource)
	if !assert.Equal(t, fooContents, string(trackedFooContents), "Mismatch on contents of foo before change") {
		return
	}

	// Check text for line.
	lines := strings.Split(fooContents, "\n")
	for index, line := range lines {
		lineText, _ := tracker.TextForLine(index, fooSource, compilercommon.SourceMapTracked)
		if !assert.Equal(t, line, lineText, "Mismatch on line text for line #%v of foo before change", index) {
			return
		}
	}

	// Check diffs.
	diffs, _ := tracker.DiffFromTracked(fooSource)
	if !assert.Equal(t, 0, len(diffs), "Expected no diffs for foo before change") {
		return
	}

	// Check offset position.
	initialRunePosition := len(fooContents) - 1
	position := fooSource.PositionForRunePosition(initialRunePosition, tracker)
	offsetPosition, _ := tracker.GetPositionOffset(position, TrackedFilePosition)
	offsetPositionRune, _ := offsetPosition.RunePosition()

	if !assert.Equal(t, offsetPositionRune, initialRunePosition, "Expected no difference in rune position for foo before change") {
		return
	}

	// Check position mapping.
	assertPositionMapping(t, tracker, fooSource, fooContents, compilercommon.SourceMapTracked)
	assertPositionMapping(t, tracker, fooSource, fooContents, compilercommon.SourceMapCurrent)

	// Change the text in the path loader.
	updatedContents := "some awesome\nfoo\nfile"
	pathLoader.setFile("foo.txt", updatedContents)

	// Check position mapping.
	assertPositionMapping(t, tracker, fooSource, fooContents, compilercommon.SourceMapTracked)
	assertPositionMapping(t, tracker, fooSource, updatedContents, compilercommon.SourceMapCurrent)

	// Check current again to ensure caching.
	assertPositionMapping(t, tracker, fooSource, updatedContents, compilercommon.SourceMapCurrent)

	// Make sure contents have not changed, but offsets and diffs have.
	trackedFooContents, _ = tracker.LoadedContents(fooSource)
	if !assert.Equal(t, fooContents, string(trackedFooContents), "Mismatch on contents of foo after change") {
		return
	}

	// Check text for line.
	for index, line := range lines {
		lineText, _ := tracker.TextForLine(index, fooSource, compilercommon.SourceMapTracked)
		if !assert.Equal(t, line, lineText, "Mismatch on line text for line #%v of foo after change", index) {
			return
		}
	}

	// Check diffs.
	diffs, _ = tracker.DiffFromTracked(fooSource)
	if !assert.True(t, len(diffs) > 0, "Expected diffs for foo after change") {
		return
	}

	// Check position offset.
	offsetPosition, _ = tracker.GetPositionOffset(position, TrackedFilePosition)
	offsetPositionRune, _ = offsetPosition.RunePosition()

	if !assert.Equal(t, offsetPositionRune, initialRunePosition+len(" awesome"), "Mismatch in computed offset position after change") {
		return
	}

	// Change the text in the path loader.
	pathLoader.setFile("foo.txt", "foo\nfile")

	// Check position offset.
	offsetPosition, _ = tracker.GetPositionOffset(position, TrackedFilePosition)
	offsetPositionRune, _ = offsetPosition.RunePosition()

	if !assert.Equal(t, offsetPositionRune, initialRunePosition-len("some\n"), "Mismatch in computed offset position after change") {
		return
	}

	// Check with a current file position.
	reversePosition, _ := tracker.GetPositionOffset(offsetPosition, CurrentFilePosition)
	reversePositionRune, _ := reversePosition.RunePosition()

	if !assert.Equal(t, reversePositionRune, initialRunePosition, "Mismatch in reverse computed offset position after change") {
		return
	}
}
