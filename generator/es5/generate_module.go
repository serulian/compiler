// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/typegraph"

	"github.com/cevaris/ordered_map"
)

// generateModules generates all the given modules into ES5.
func (gen *es5generator) generateModules(modules []typegraph.TGModule) map[typegraph.TGModule]string {
	// Queue all the modules to be generated.
	generatedSource := make([]string, len(modules))
	queue := compilerutil.Queue()
	for index, module := range modules {
		fn := func(key interface{}, value interface{}) bool {
			generatedSource[key.(int)] = gen.generateModule(value.(typegraph.TGModule))
			return true
		}

		queue.Enqueue(index, module, fn)
	}

	// Generate the full source tree for each module.
	queue.Run()

	// Build a map from module to source tree.
	moduleMap := map[typegraph.TGModule]string{}
	for index, module := range modules {
		moduleMap[module] = generatedSource[index]
	}

	return moduleMap
}

// generateModule generates the given module into ES5.
func (gen *es5generator) generateModule(module typegraph.TGModule) string {
	generating := generatingModule{module, gen}
	return gen.templater.Execute("module", moduleTemplateStr, generating)
}

// generatingModule represents a module being generated.
type generatingModule struct {
	Module    typegraph.TGModule
	Generator *es5generator // The parent generator.
}

// ExportedPath returns the full exported path for this module.
func (gm generatingModule) ExportedPath() string {
	return gm.Generator.pather.GetRelativeModulePath(gm.Module)
}

// GenerateMembers generates the source for all the implemented members defined under the module.
func (gm generatingModule) GenerateMembers() *ordered_map.OrderedMap {
	return gm.Generator.generateImplementedMembers(gm.Module)
}

// GenerateTypes generates the source for all the types defined under the module.
func (gm generatingModule) GenerateTypes() *ordered_map.OrderedMap {
	return gm.Generator.generateTypes(gm.Module)
}

// GenerateVariables generates the source for all the variables defined under the module.
func (gm generatingModule) GenerateVariables() *ordered_map.OrderedMap {
	return gm.Generator.generateVariables(gm.Module)
}

// moduleTemplateStr defines the template for generating a module.
const moduleTemplateStr = `
$module('{{ .ExportedPath }}', function() {
  var $static = this;

  {{range $idx, $kv := .GenerateTypes.Iter }}
  	{{ $kv.Value }};
  {{end}}
  
  {{range $idx, $kv := .GenerateMembers.Iter }}
  	{{ $kv.Value }};
  {{end}}

  {{range $idx, $kv := .GenerateVariables.Iter }}
  	this.$init(function() {
		return ({{ $kv.Value }});
	});
  {{end}}
});
`
