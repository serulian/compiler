// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"strings"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// FindCommentedNode attempts to find the node in the SRG with the given comment attached.
func (g *SRG) FindCommentedNode(commentValue string) (compilergraph.GraphNode, bool) {
	return g.layer.StartQuery(commentValue).
		In(parser.NodeCommentPredicateValue).
		In(parser.NodePredicateChild).
		TryGetNode()
}

// FindComment attempts to find the comment in the SRG with the given value.
func (g *SRG) FindComment(commentValue string) (SRGComment, bool) {
	commentNode, ok := g.layer.StartQuery(commentValue).
		In(parser.NodeCommentPredicateValue).
		TryGetNode()

	if !ok {
		return SRGComment{}, false
	}

	return SRGComment{commentNode, g}, true
}

// AllComments returns an iterator over all the comment nodes found in the SRG.
func (g *SRG) AllComments() compilergraph.NodeIterator {
	return g.layer.StartQuery().Has(parser.NodeCommentPredicateValue).BuildNodeIterator()
}

// getDocumentationCommentForNode returns the last comment attached to the given SRG node, if any.
func (g *SRG) getDocumentationCommentForNode(node compilergraph.GraphNode) (SRGComment, bool) {
	// The documentation on a node is its *last* comment.
	cit := node.StartQuery().
		Out(parser.NodePredicateChild).
		Has(parser.NodeCommentPredicateValue).
		BuildNodeIterator()

	var comment SRGComment
	var commentFound = false
	for cit.Next() {
		commentFound = true
		comment = SRGComment{cit.Node(), g}
	}

	return comment, commentFound
}

// srgCommentKind defines a struct for specifying various pieces of config about a kind of
// comment found in the SRG.
type srgCommentKind struct {
	prefix          string
	suffix          string
	multilinePrefix string
	isDocComment    bool
}

// commentsKinds defines the various kinds of comments.
var commentsKinds = []srgCommentKind{
	srgCommentKind{"/**", "*/", "*", true},
	srgCommentKind{"/*", "*/", "*", false},
	srgCommentKind{"//", "", "", false},
}

// getCommentKind returns the kind of the comment, if any.
func getCommentKind(value string) (srgCommentKind, bool) {
	for _, kind := range commentsKinds {
		if strings.HasPrefix(value, kind.prefix) {
			return kind, true
		}
	}

	return srgCommentKind{}, false
}

// getTrimmedCommentContentsString returns trimmed contents of the given comment.
func getTrimmedCommentContentsString(value string) string {
	kind, isValid := getCommentKind(value)
	if !isValid {
		return ""
	}

	lines := strings.Split(value, "\n")

	var newLines = make([]string, 0, len(lines))
	for index, line := range lines {
		var updatedLine = line
		if index == 0 {
			updatedLine = strings.TrimPrefix(strings.TrimSpace(updatedLine), kind.prefix)
		}

		if index == len(lines)-1 {
			updatedLine = strings.TrimSuffix(strings.TrimSpace(updatedLine), kind.suffix)
		}

		if index > 0 && index < len(lines)-1 {
			updatedLine = strings.TrimPrefix(strings.TrimSpace(updatedLine), kind.multilinePrefix)
		}

		updatedLine = strings.TrimSpace(updatedLine)
		newLines = append(newLines, updatedLine)
	}

	return strings.TrimSpace(strings.Join(newLines, "\n"))
}

// SRGComment wraps a comment node with helpers for accessing documentation.
type SRGComment struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// kind returns the kind of the comment.
func (c SRGComment) kind() srgCommentKind {
	value := c.GraphNode.Get(parser.NodeCommentPredicateValue)
	kind, isValid := getCommentKind(value)
	if !isValid {
		panic("Invalid comment kind")
	}
	return kind
}

// IsDocComment returns true if the comment is a doc comment, instead of a normal
// comment.
func (c SRGComment) IsDocComment() bool {
	return c.kind().isDocComment
}

// ParentNode returns the node that contains this comment.
func (c SRGComment) ParentNode() compilergraph.GraphNode {
	return c.GraphNode.GetIncomingNode(parser.NodePredicateChild)
}

// Contents returns the trimmed contents of the comment.
func (c SRGComment) Contents() string {
	value := c.GraphNode.Get(parser.NodeCommentPredicateValue)
	return getTrimmedCommentContentsString(value)
}

// Documentation returns the documentation found in this comment, if any.
func (c SRGComment) Documentation() (SRGDocumentation, bool) {
	return SRGDocumentation{
		parentComment: c,
		commentValue:  c.Contents(),
		isDocComment:  c.IsDocComment(),
	}, true
}

// SRGDocumentation represents documentation found on an SRG node. It is distinct from a comment in
// that it can be *part* of another comment.
type SRGDocumentation struct {
	parentComment SRGComment
	commentValue  string
	isDocComment  bool
}

// String returns the human-readable documentation string.
func (d SRGDocumentation) String() string {
	return d.commentValue
}

// IsDocComment returns true if the documentation was found in a doc-comment style
// comment (instead of a standard comment).
func (d SRGDocumentation) IsDocComment() bool {
	return d.isDocComment
}
