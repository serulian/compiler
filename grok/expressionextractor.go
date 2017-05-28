// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"strings"
)

// extractorState represents the various states of the expression extractor.
type extractorState int

const (
	extractorStateNormal extractorState = iota
	extractorStateInString
)

// nestingRunes defines all runes that perform nesting, along with their corresponding
// unnested runes.
var nestingRunes = map[rune]rune{
	'{': '}',
	'[': ']',
	'(': ')',
}

// stringDelimeters defines all runes that are string delimeters.
var stringDelimeters = map[rune]bool{
	'"':  true,
	'\'': true,
	'`':  true,
}

// expressionExtractor is a helper type for extracting expressions from an input string
// that may represent a partial expression.
type expressionExtractor struct {
	inputString       string
	currentState      extractorState
	nestingStack      *nestingStack
	currentStringRune rune
}

// extractExpression extracts the closest inner expression from the given input string, if any.
// For example, this method will perform the following transformations:
// a -> a
// a.b ->
// a[1] -> a[1]
// a[b -> b (partial)
// a(b -> b (partial)
// a(b] -> *error*
// if something -> something (partial after the space)
func extractExpression(inputString string) (string, bool) {
	expressionExtractor := &expressionExtractor{
		inputString:       strings.TrimSpace(inputString),
		currentState:      extractorStateNormal,
		nestingStack:      &nestingStack{},
		currentStringRune: ' ',
	}
	return expressionExtractor.extract()
}

func (ee *expressionExtractor) extract() (string, bool) {
	closingRunes := map[rune]bool{}
	for _, closeRune := range nestingRunes {
		closingRunes[closeRune] = true
	}

	afterSpaceStartIndex := 0
	for index, char := range ee.inputString {
		_, isOpenRune := nestingRunes[char]
		_, isCloseRune := closingRunes[char]
		_, isStringRune := stringDelimeters[char]

		switch {
		// If a string rune, then we are either entering or exiting a string (unless escaped).
		case isStringRune:
			// If in a normal state, then we are entering a string.
			if ee.currentState == extractorStateNormal {
				ee.currentStringRune = char
				ee.currentState = extractorStateInString
			} else {
				if index > 0 {
					// Check for an escape character.
					if ee.inputString[index-1] != '\\' {
						// Ensure we have a matching rune.
						if char != ee.currentStringRune {
							// String mismatch.
							return "", false
						}

						ee.currentState = extractorStateNormal
					}
				}
			}

		case ee.currentState == extractorStateInString:
			// If inside a string, just consume.
			break

		case isOpenRune:
			ee.nestingStack.push(char, index)

		case isCloseRune:
			openRune, _ := ee.nestingStack.pop()
			matchingCloseRune, _ := nestingRunes[openRune]
			if matchingCloseRune != char {
				// Nesting rune mismatch.
				return "", false
			}

		case char == ' ':
			if ee.currentState == extractorStateNormal {
				afterSpaceStartIndex = index + 1
			}
		}
	}

	if ee.currentState != extractorStateNormal {
		return "", false
	}

	// If we've reached this point, return the expression string found
	// at the last nesting or the last space index.
	_, lastNestingStartIndex, _ := ee.nestingStack.topValue()
	expressionStartIndex := lastNestingStartIndex + 1
	if afterSpaceStartIndex > expressionStartIndex {
		expressionStartIndex = afterSpaceStartIndex
	}

	return ee.inputString[expressionStartIndex:], true
}
