// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/generator/escommon/esbuilder"
	"github.com/serulian/compiler/graphs/typegraph"

	"github.com/cevaris/ordered_map"
)

var _ = fmt.Printf

// generateModules generates all the given modules into ES5.
func (gen *es5generator) generateModules(modules []typegraph.TGModule) map[typegraph.TGModule]esbuilder.SourceBuilder {
	// Queue all the modules to be generated.
	generatedSource := make([]esbuilder.SourceBuilder, len(modules))
	queue := compilerutil.Queue()
	for index, module := range modules {
		fn := func(key interface{}, value interface{}, cancel compilerutil.CancelFunction) bool {
			generatedSource[key.(int)] = gen.generateModule(value.(typegraph.TGModule))
			return true
		}

		queue.Enqueue(index, module, fn)
	}

	// Generate the full source tree for each module.
	queue.Run()

	// Build a map from module to source tree.
	moduleMap := map[typegraph.TGModule]esbuilder.SourceBuilder{}
	for index, module := range modules {
		moduleMap[module] = generatedSource[index]
	}

	return moduleMap
}

// generateModule generates the given module into ES5.
func (gen *es5generator) generateModule(module typegraph.TGModule) esbuilder.SourceBuilder {
	generating := generatingModule{module, gen}
	return esbuilder.Template("module", moduleTemplateStr, generating)
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
	orderedMap := ordered_map.NewOrderedMap()
	gm.Generator.generateVariables(gm.Module, orderedMap)
	return orderedMap
}

// recursivelyCollectInitDependencies recursively collects the initialization dependency fields for the
// field member, placing them in the deps map.
func (gm generatingModule) recursivelyCollectInitDependencies(current typegraph.TGMember, field typegraph.TGMember, deps map[typegraph.TGMember]bool) {
	if _, found := deps[current]; found {
		return
	}

	if field.NodeId != current.NodeId && current.IsField() {
		deps[current] = true
	} else {
		srgMember, hasSRGMember := gm.Generator.getSRGMember(current)
		if !hasSRGMember {
			return
		}

		scope, _ := gm.Generator.scopegraph.GetScope(srgMember.GraphNode)
		for _, staticDep := range scope.GetStaticDependencies() {
			memberNodeId := compilergraph.GraphNodeId(staticDep.GetReferencedNode())
			member := gm.Generator.scopegraph.TypeGraph().GetTypeOrMember(memberNodeId)
			gm.recursivelyCollectInitDependencies(member.(typegraph.TGMember), field, deps)
		}
	}
}

// FieldId returns a stable, unique ID for the given field.
func (gm generatingModule) FieldId(member typegraph.TGMember) string {
	srgMember, _ := gm.Generator.getSRGMember(member)
	return srgMember.UniqueId()
}

// InitDependenceis returns a slice containing the FieldId's of the fields that the given
// field depends upon for initialization, if any.
func (gm generatingModule) InitDependencies(field typegraph.TGMember) []string {
	deps := map[typegraph.TGMember]bool{}
	gm.recursivelyCollectInitDependencies(field, field, deps)

	dependencies := make([]string, 0, len(deps))
	for mem, _ := range deps {
		srgMember, hasSRGMember := gm.Generator.getSRGMember(mem)
		if hasSRGMember && srgMember.IsStatic() {
			dependencies = append(dependencies, srgMember.UniqueId())
		}
	}

	return dependencies
}

// moduleTemplateStr defines the template for generating a module.
const moduleTemplateStr = `
{{ $types := .GenerateTypes }}
{{ $members := .GenerateMembers }}
{{ $vars := .GenerateVariables }}

{{ $hasContents := or $types.Len $members.Len $vars.Len }}

{{ if $hasContents }}
$module('{{ .ExportedPath }}', function() {
  var $static = this;

  {{range $idx, $kv := $types.UnsafeIter }}
  	{{ emit $kv.Value }};
  {{end}}
  
  {{range $idx, $kv := $members.UnsafeIter }}
  	{{ emit $kv.Value }};
  {{end}}

  {{ $parent := . }}

  {{range $idx, $kv := $vars.UnsafeIter }}
  	this.$init(function() {
  		{{ if $kv.Value.IsAsync }}
			return ({{ emit $kv.Value.Source }});
  		{{ else }}
  		    return $promise.new(function (resolve) {
				{{ emit $kv.Value.Source }};
				resolve();
		    });
  		{{ end }}
	}, '{{ $parent.FieldId $kv.Key }}', [{{ range $ddx, $did := $parent.InitDependencies $kv.Key }}{{ if $ddx }}, {{ end }}'{{ $did }}'{{ end }}]);
  {{end}}
});
{{ end }}
`
