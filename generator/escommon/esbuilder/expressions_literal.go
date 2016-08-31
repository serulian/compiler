// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esbuilder

import (
	"fmt"
	"strconv"
)

// literalNode defines a literal value.
type literalNode struct {
	// value is the literal value.
	value string
}

// arrayNode defines an array literal value.
type arrayNode struct {
	// values are the values in the array.
	values []ExpressionBuilder
}

// objectNode defines an object literal value.
type objectNode struct {
	// entries are the entries in the object.
	entries []ExpressionBuilder
}

// objectNodeEntry defines a single entry in the object.
type objectNodeEntry struct {
	// key is the key for the object entry.
	key string

	// value is the value for the object entry.
	value ExpressionBuilder
}

func (node literalNode) emit(sb *sourceBuilder) {
	sb.append(node.value)
}

func (node arrayNode) emit(sb *sourceBuilder) {
	sb.append("[")
	sb.emitSeparated(node.values, ",")
	sb.append("]")
}

func (node objectNode) emit(sb *sourceBuilder) {
	sb.append("{")
	sb.emitSeparated(node.entries, ",")
	sb.append("}")
}

func (node objectNodeEntry) emit(sb *sourceBuilder) {
	sb.emit(Value(node.key))
	sb.append(":")
	sb.emitWrapped(node.value)
}

func (node literalNode) isStateless() bool {
	return true
}

func (node arrayNode) isStateless() bool {
	for _, expr := range node.values {
		if !expr.IsStateless() {
			return false
		}
	}
	return true
}

func (node objectNode) isStateless() bool {
	for _, entry := range node.entries {
		if !entry.IsStateless() {
			return false
		}
	}
	return true
}

func (node objectNodeEntry) isStateless() bool {
	return node.value.IsStateless()
}

// ObjectEntry returns an object literal entry.
func ObjectEntry(key string, value ExpressionBuilder) ExpressionBuilder {
	return expressionBuilder{objectNodeEntry{key, value}, nil}
}

// Object returns an object literal.
func Object(entries ...ExpressionBuilder) ExpressionBuilder {
	return expressionBuilder{objectNode{entries}, nil}
}

// Array returns an array literal.
func Array(values ...ExpressionBuilder) ExpressionBuilder {
	return expressionBuilder{arrayNode{values}, nil}
}

// LiteralValue returns a literal value.
func LiteralValue(value string) ExpressionBuilder {
	return expressionBuilder{literalNode{value}, nil}
}

// Value returns a literal value.
func Value(value interface{}) ExpressionBuilder {
	switch t := value.(type) {
	case bool:
		if t {
			return LiteralValue("true")
		} else {
			return LiteralValue("false")
		}

	case int:
		return LiteralValue(strconv.Itoa(t))

	case float64:
		return LiteralValue(strconv.FormatFloat(t, 'E', -1, 32))

	case string:
		return LiteralValue(strconv.Quote(t))

	default:
		panic(fmt.Sprintf("unexpected value type %T\n", t))
	}
}
