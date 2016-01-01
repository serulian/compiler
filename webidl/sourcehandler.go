// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webidl

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/webidl/parser"
)

// irgSourceHandler implements the SourceHandler interface from the packageloader for
// populating the WebIDL IRG from webidl files.
type irgSourceHandler struct {
	irg *WebIRG // The IRG being populated.
}

func (sh *irgSourceHandler) Kind() string {
	return "webidl"
}

func (sh *irgSourceHandler) PackageFileExtension() string {
	return ".webidl"
}

func (sh *irgSourceHandler) Parse(source compilercommon.InputSource, input string, importHandler packageloader.ImportHandler) {
	parser.Parse(sh.irg.buildASTNode, source, input)
}

func (sh *irgSourceHandler) Verify(packageMap map[string]packageloader.PackageInfo, errorReporter packageloader.ErrorReporter, warningReporter packageloader.WarningReporter) {

}
