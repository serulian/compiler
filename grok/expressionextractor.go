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
//
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

// extractCalled extracts the closest *called* inner expression from the given input string, if any.
//
// For example, this method will perform the following transformations:
// a -> (none)
// a.b -> (none)
// a[1] -> (none)
// a[b -> a (w/index 0)
// a(b -> a (w/index 0)
// a(b, c -> a (w/index 1)
// a( -> a (w/index 0)
func extractCalled(inputString string) (string, int, bool) {
	expressionExtractor := &expressionExtractor{
		inputString:       strings.TrimSpace(inputString),
		currentState:      extractorStateNormal,
		nestingStack:      &nestingStack{},
		currentStringRune: ' ',
	}

	// Extract the normal expression.
	_, ok := expressionExtractor.extract()
	if !ok {
		return "", -1, false
	}

	// Find the containing nesting level, to determine it expression.
	containingNesting := expressionExtractor.nestingStack.top

	// Pop it.
	expressionExtractor.nestingStack.pop()

	// Find the second nesting level, if any.
	if expressionExtractor.nestingStack.top == nil {
		return "", -1, false
	}

	// The called expression is that found between the two levels.
	startIndex := expressionExtractor.nestingStack.top.expressionStartIndex()
	endIndex := containingNesting.startIndex
	commaIndex := containingNesting.commaCounter
	return inputString[startIndex:endIndex], commaIndex, true
}

func (ee *expressionExtractor) extract() (string, bool) {
	// Build the map of closing runes.
	closingRunes := map[rune]bool{}
	for _, closeRune := range nestingRunes {
		closingRunes[closeRune] = true
	}

	// Push a top-level "nesting" state onto the stack.
	ee.nestingStack.push(' ', -1)

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

		case char == ' ' || char == ',':
			if ee.currentState == extractorStateNormal {
				ee.nestingStack.top.lastSeparatorIndex = index
				if char == ',' {
					ee.nestingStack.top.commaCounter++
				}
			}
		}
	}

	if ee.currentState != extractorStateNormal {
		return "", false
	}

	// If we've reached this point, return the expression string found
	// at the last nesting or the last separator index.
	expressionStartIndex := ee.nestingStack.top.expressionStartIndex()
	return ee.inputString[expressionStartIndex:], true
}
