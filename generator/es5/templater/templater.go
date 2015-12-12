// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The templater package implements helper methods for working with go templates.
package templater

import (
	"bytes"
	"text/template"
)

// Templater is a helper type for running go templates.
type Templater struct {
	templateCache map[string]*template.Template
}

func New() *Templater {
	return &Templater{
		templateCache: map[string]*template.Template{},
	}
}

// execute runs the given go-template over the given context.
func (temp *Templater) Execute(name string, templateStr string, context interface{}) string {
	var parsed *template.Template = nil

	// Check the cache.
	cached, isCached := temp.templateCache[templateStr]
	if !isCached {
		t := template.New(name)
		parsedTemplate, err := t.Parse(templateStr)
		if err != nil {
			panic(err)
		}

		temp.templateCache[templateStr] = parsedTemplate
		parsed = parsedTemplate
	} else {
		parsed = cached
	}

	// Execute the template.
	var source bytes.Buffer
	eerr := parsed.Execute(&source, context)
	if eerr != nil {
		panic(eerr)
	}

	return source.String()
}
