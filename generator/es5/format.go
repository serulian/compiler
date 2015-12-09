// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"os"

	"github.com/robertkrimen/otto"
)

var ottovm *otto.Otto // The otto VM.

func formatSource(source string) string {
	if ottovm == nil {
		ottovm = otto.New()
		file, err := os.Open("jsbeautifier.js")
		defer file.Close()

		if err != nil {
			panic(err)
		}

		_, err2 := ottovm.Run(`exports = {}`)
		if err2 != nil {
			panic(err2)
		}

		_, err3 := ottovm.Run(file)
		if err != nil {
			panic(err3)
		}
	}

	ottovm.Set("sourceCode", source)
	formatted, err := ottovm.Run(`exports.js_beautify(sourceCode, {
		'indent_size': 2,
		'preserve_newlines': false,
		'jslint_happy': true,
		'wrap_line_length': 120,
		'end_with_newline': true
	})`)

	if err != nil {
		panic(err)
	}

	return formatted.String()
}
