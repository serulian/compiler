// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package developer

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/serulian/compiler/builder"
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/generator/es5"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/sourcemap"
)

// developerRuntimeTemplateStr defines the template for the ES code emitted to allow for
// UI when running the developer tool.
const developerRuntimeTemplateStr = `
	console.group('Compilation of project {{ .Name }}');
	console.info('Compilation has begun');
	document.write('<script src="http://localhost{{ .Addr }}/{{ .Name }}.develop.js"></script>');
`

// developTransaction represents a single transaction of loading source via the development
// server.
type developTransaction struct {
	vcsDevelopmentDirectories []string             // VCS development directory paths.
	rootSourceFilePath        string               // The root source file
	addr                      string               // The address of the running server.
	name                      string               // The name of the source being developed.
	offsetCount               int                  // The number of emitted call lines that offsets the generated source.
	sourceMap                 *sourcemap.SourceMap // The constructed source map.
}

func newDevelopTransaction(rootSourceFilePath string, vcsDevelopmentDirectories []string, addr string, name string) *developTransaction {
	return &developTransaction{
		vcsDevelopmentDirectories: vcsDevelopmentDirectories,
		rootSourceFilePath:        rootSourceFilePath,
		addr:                      addr,
		name:                      name,
		offsetCount:               0,
		sourceMap:                 nil,
	}
}

// Start starts the develop transaction by emitting the runtime that logs the build
// and requests the built data URL (which actually starts the build).
func (dt *developTransaction) Start(w http.ResponseWriter, r *http.Request) {
	// Emit the runtime.
	context := struct {
		Addr string
		Name string
	}{dt.addr, dt.name}

	// Write the runtime template.
	t := template.New("runtimeTemplate")
	parsedTemplate, err := t.Parse(developerRuntimeTemplateStr)
	if err != nil {
		panic(err)
	}

	var source bytes.Buffer
	eerr := parsedTemplate.Execute(&source, context)
	if eerr != nil {
		panic(eerr)
	}

	fmt.Fprint(w, source.String())
}

// Build performs the build of the source, writing the result to the response writer.
func (dt *developTransaction) Build(w http.ResponseWriter, r *http.Request) {
	// Build a scope graph for the project. This will conduct parsing and type graph
	// construction on our behalf.
	scopeResult, err := scopegraph.ParseAndBuildScopeGraph(dt.rootSourceFilePath,
		dt.vcsDevelopmentDirectories,
		builder.CORE_LIBRARY)

	if err != nil {
		dt.emitInfo(w, "Build failed: %s", err)
		dt.closeGroup(w)
	}

	if !scopeResult.Status {
		dt.sourceMap = sourcemap.NewSourceMap(dt.name+".develop.js", "source/")

		for _, warning := range scopeResult.Warnings {
			dt.emitWarning(w, warning)
		}

		for _, err := range scopeResult.Errors {
			dt.emitError(w, err)
		}

		dt.emitInfo(w, "Build failed")
		dt.closeGroup(w)
	} else {
		// Generate the program's source.
		generated, sourceMap, err := es5.GenerateES5(scopeResult.Graph, dt.name+".develop.js", "source/")
		if err != nil {
			panic(err)
		}

		dt.sourceMap = sourceMap

		fmt.Fprint(w, generated)
		dt.emitInfo(w, "Build completed successfully")
		dt.closeGroup(w)

		dt.offsetCount = len(strings.Split(string(generated), "\n")) + 1

		for _, warning := range scopeResult.Warnings {
			dt.emitWarning(w, warning)
		}
	}

	fmt.Fprintf(w, "//# sourceMappingURL=/%s.develop.js.map\n", dt.name)
}

func (dt *developTransaction) ServeSourceMap(w http.ResponseWriter, r *http.Request) {
	marshalled, err := dt.sourceMap.Build().Marshal()
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(marshalled))
}

func (dt *developTransaction) ServeSourceFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	path := params["path"]

	fullPath := filepath.Join(filepath.Dir(dt.rootSourceFilePath), string(path))
	if _, err := os.Stat(fullPath); err != nil {
		fullPath = "/" + path
	}

	contents, err := ioutil.ReadFile(fullPath)
	if err != nil {
		fmt.Fprintf(w, "Error when trying to read source file %s: %v", fullPath, err)
	}

	fmt.Fprint(w, string(contents))
}

func (dt *developTransaction) emitError(w http.ResponseWriter, sourceErr compilercommon.SourceError) {
	message := fmt.Sprintf("console.error(%v);\n", strconv.Quote(sourceErr.Error()))
	fmt.Fprint(w, message)

	line, col, err := sourceErr.SourceRange().Start().LineAndColumn()
	if err == nil {
		dt.sourceMap.AddMapping(dt.offsetCount, 0, sourcemap.SourceMapping{
			SourcePath:     string(sourceErr.SourceRange().Source()),
			LineNumber:     line,
			ColumnPosition: col,
		})

		dt.offsetCount++
	}
}

func (dt *developTransaction) emitWarning(w http.ResponseWriter, warn compilercommon.SourceWarning) {
	message := fmt.Sprintf("console.warn(%v);\n", strconv.Quote(warn.String()))
	fmt.Fprint(w, message)

	line, col, err := warn.SourceRange().Start().LineAndColumn()
	if err == nil {
		dt.sourceMap.AddMapping(dt.offsetCount, 0, sourcemap.SourceMapping{
			SourcePath:     string(warn.SourceRange().Source()),
			LineNumber:     line,
			ColumnPosition: col,
		})

		dt.offsetCount++
	}
}

func (dt *developTransaction) emitInfo(w http.ResponseWriter, msg string, args ...interface{}) {
	fmt.Fprintf(w, "console.info('%v');\n", fmt.Sprintf(msg, args...))
	dt.offsetCount++
}

func (dt *developTransaction) closeGroup(w http.ResponseWriter) {
	fmt.Fprintf(w, "console.groupEnd();\n")
	dt.offsetCount++
}
