// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"

	"strings"

	"github.com/fatih/color"
	"github.com/kr/text"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

type ConsoleLogLevel int

const INDENTATION = "  "

const (
	InfoLogLevel ConsoleLogLevel = iota
	SuccessLogLevel
	WarningLogLevel
	ErrorLogLevel
)

// LogToConsole logs a message to the console, with the given level and range.
func LogToConsole(level ConsoleLogLevel, sourceRange compilercommon.SourceRange, message string, args ...interface{}) {
	width, _ := terminal.Width()
	if width <= 0 {
		width = 80
	}

	locationColor := color.New(color.FgWhite)
	messageColor := color.New(color.FgHiWhite)
	codeColor := color.New(color.FgWhite)
	prefixColor := color.New(color.FgWhite)

	prefixText := ""

	switch level {
	case InfoLogLevel:
		prefixColor = color.New(color.FgBlue, color.Bold)
		prefixText = "INFO: "

	case SuccessLogLevel:
		prefixColor = color.New(color.FgGreen, color.Bold)
		prefixText = "SUCCESS: "

	case WarningLogLevel:
		prefixColor = color.New(color.FgYellow, color.Bold)
		prefixText = "WARNING: "

	case ErrorLogLevel:
		prefixColor = color.New(color.FgRed, color.Bold)
		prefixText = "ERROR: "
	}

	formattedMessage := fmt.Sprintf(message, args...)

	if sourceRange == nil {
		prefixColor.Print(prefixText)
		messageColor.Printf("%s\n", formattedMessage)
		return
	}

	startLine, startCol, err := sourceRange.Start().LineAndColumn()
	if err != nil {
		startLine = 0
		startCol = 0
	}

	endLine, endCol, err := sourceRange.End().LineAndColumn()
	if err != nil {
		endLine = 0
		endCol = 0
	}

	locationText := fmt.Sprintf("%v:%v:%v:", sourceRange.Source(), startLine+1, startCol+1)

	// If the text will go past the terminal width, then make it multiline and add extra whitespace
	// after it.
	if len(prefixText+locationText+formattedMessage) > int(width) {
		prefixColor.Print(prefixText)
		locationColor.Print(locationText)
		messageColor.Printf("\n%s\n\n", text.Indent(text.Wrap(formattedMessage, int(width)-len(INDENTATION)), INDENTATION))
	} else {
		prefixColor.Print(prefixText)
		locationColor.Print(locationText)
		messageColor.Printf(" %s\n", formattedMessage)
	}

	// Print the actual code line, and annotate the range.
	lineContents, err := sourceRange.Start().LineText()
	if endLine != startLine {
		endCol = len(lineContents)
	}

	if err == nil {
		caretValue := ""
		for index := range lineContents {
			if index < startCol {
				caretValue += " "
			} else if index >= startCol && index <= endCol {
				caretValue += "^"
			}
		}

		codeColor.Printf("%s\n%s\n\n\n", text.Indent(strings.Replace(lineContents, "\t", " ", -1), INDENTATION), text.Indent(caretValue, INDENTATION))
	} else {
		codeColor.Printf("\n\n")
	}
}
