// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"

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

// LogToConsole logs a message to the console, with the given level and SourceAndLocation.
func LogToConsole(level ConsoleLogLevel, sal compilercommon.SourceAndLocation, message string, args ...interface{}) {
	width, _ := terminal.Width()
	if width <= 0 {
		width = 80
	}

	locationColor := color.New(color.FgWhite)
	messageColor := color.New(color.FgHiWhite)

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

	locationText := fmt.Sprintf("At %v:%v:%v:", sal.Source(), sal.Location().LineNumber()+1, sal.Location().ColumnPosition()+1)
	formattedMessage := fmt.Sprintf(message, args...)

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
}
