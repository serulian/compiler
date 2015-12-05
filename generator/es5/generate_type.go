// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/typegraph"
)

// generateTypes generates all the types under the  given modules into ES5.
func (gen *es5generator) generateTypes(module typegraph.TGModule) map[typegraph.TGTypeDecl]string {
	// Queue all the types to be generated.
	types := module.Types()
	generatedSource := make([]string, len(types))
	queue := compilerutil.Queue()
	for index, typedef := range types {
		fn := func(key interface{}, value interface{}) bool {
			generatedSource[key.(int)] = gen.generateType(value.(typegraph.TGTypeDecl))
			return true
		}

		queue.Enqueue(index, typedef, fn)
	}

	// Generate the full source tree for each type.
	queue.Run()

	// Build a map from type to source tree.
	typeMap := map[typegraph.TGTypeDecl]string{}
	for index, typedef := range types {
		typeMap[typedef] = generatedSource[index]
	}

	return typeMap
}

// generateType generates the given type into ES5.
func (gen *es5generator) generateType(typedef typegraph.TGTypeDecl) string {
	generating := generatingType{typedef, gen}

	switch typedef.TypeKind() {
	case typegraph.ClassType:
		return gen.runTemplate("class", classTemplateStr, generating)

	case typegraph.InterfaceType:
		return gen.runTemplate("interface", interfaceTemplateStr, generating)

	default:
		panic("Unknown typedef kind")
		return ""
	}
}

// generatingType represents a type being generated.
type generatingType struct {
	Type      typegraph.TGTypeDecl
	Generator *es5generator // The parent generator.
}

// classTemplateStr defines the template for generating a class type.
const classTemplateStr = `
$parent.cls('{{ .Context.Type.Name }}', function() {
});
`

// interfaceTemplateStr defines the template for generating an interface type.
const interfaceTemplateStr = `
$parent.interface('{{ .Context.Type.Name }}', function() {
});
`
