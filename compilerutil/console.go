// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"

	"strings"

	color "github.com/fatih/color"
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

var (
	InfoColor       = color.New(color.FgBlue, color.Bold)
	SuccessColor    = color.New(color.FgGreen, color.Bold)
	WarningColor    = color.New(color.FgYellow, color.Bold)
	ErrorColor      = color.New(color.FgRed, color.Bold)
	BoldWhiteColor  = color.New(color.FgWhite, color.Bold)
	FaintWhiteColor = color.New(color.FgWhite, color.Faint)
	MessageColor    = color.New(color.FgHiWhite)
)

// LogToConsole logs a message to the console, with the given level and range.
func LogToConsole(level ConsoleLogLevel, sourceRange compilercommon.SourceRange, message string, args ...interface{}) {
	prefixColor := color.New(color.FgWhite)
	prefixText := ""

	switch level {
	case InfoLogLevel:
		prefixColor = InfoColor
		prefixText = "INFO: "

	case SuccessLogLevel:
		prefixColor = SuccessColor
		prefixText = "SUCCESS: "

	case WarningLogLevel:
		prefixColor = WarningColor
		prefixText = "WARNING: "

	case ErrorLogLevel:
		prefixColor = ErrorColor
		prefixText = "ERROR: "
	}

	LogMessageToConsole(prefixColor, prefixText, sourceRange, message, args...)
}

// LogMessageToConsole logs a message to the console, with the given level and range.
func LogMessageToConsole(prefixColor *color.Color, prefixText string, sourceRange compilercommon.SourceRange, message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)

	if sourceRange == nil {
		prefixColor.Print(prefixText)
		MessageColor.Printf("%s\n", formattedMessage)
		return
	}

	locationColor := color.New(color.FgWhite)
	codeColor := color.New(color.FgWhite)

	width, _ := terminal.Width()
	if width <= 0 {
		width = 80
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
		MessageColor.Printf("\n%s\n\n", text.Indent(text.Wrap(formattedMessage, int(width)-len(INDENTATION)), INDENTATION))
	} else {
		prefixColor.Print(prefixText)
		locationColor.Print(locationText)
		MessageColor.Printf(" %s\n", formattedMessage)
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
