// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"sort"
	"strings"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/sourceshape"
)

// FindCommentedNode attempts to find the node in the SRG with the given comment attached.
func (g *SRG) FindCommentedNode(commentValue string) (compilergraph.GraphNode, bool) {
	return g.layer.StartQuery(commentValue).
		In(sourceshape.NodeCommentPredicateValue).
		In(sourceshape.NodePredicateChild).
		TryGetNode()
}

// FindComment attempts to find the comment in the SRG with the given value.
func (g *SRG) FindComment(commentValue string) (SRGComment, bool) {
	commentNode, ok := g.layer.StartQuery(commentValue).
		In(sourceshape.NodeCommentPredicateValue).
		TryGetNode()

	if !ok {
		return SRGComment{}, false
	}

	return SRGComment{commentNode, g}, true
}

// AllComments returns an iterator over all the comment nodes found in the SRG.
func (g *SRG) AllComments() compilergraph.NodeIterator {
	return g.layer.StartQuery().Has(sourceshape.NodeCommentPredicateValue).BuildNodeIterator()
}

// getDocumentationForNode returns the documentation attached to the given node, if any.
func (g *SRG) getDocumentationForNode(node compilergraph.GraphNode) (SRGDocumentation, bool) {
	cit := node.StartQuery().
		Out(sourceshape.NodePredicateChild).
		Has(sourceshape.NodeCommentPredicateValue).
		BuildNodeIterator()

	var docComment *SRGComment
	for cit.Next() {
		commentNode := cit.Node()
		srgComment := SRGComment{commentNode, g}
		if docComment == nil {
			docComment = &srgComment
		}

		// Doc comments always win.
		if srgComment.IsDocComment() {
			docComment = &srgComment
		}
	}

	if docComment == nil {
		return SRGDocumentation{}, false
	}

	return docComment.Documentation()
}

// getDocumentationCommentForNode returns the last comment attached to the given SRG node, if any.
func (g *SRG) getDocumentationCommentForNode(node compilergraph.GraphNode) (SRGComment, bool) {
	// The documentation on a node is its *last* comment.
	cit := node.StartQuery().
		Out(sourceshape.NodePredicateChild).
		Has(sourceshape.NodeCommentPredicateValue).
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
	value := c.GraphNode.Get(sourceshape.NodeCommentPredicateValue)
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
	return c.GraphNode.GetIncomingNode(sourceshape.NodePredicateChild)
}

// Contents returns the trimmed contents of the comment.
func (c SRGComment) Contents() string {
	value := c.GraphNode.Get(sourceshape.NodeCommentPredicateValue)
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

type byLocation struct {
	values   []string
	strValue string
}

func (s byLocation) Len() int {
	return len(s.values)
}

func (s byLocation) Swap(i, j int) {
	s.values[i], s.values[j] = s.values[j], s.values[i]
}

func (s byLocation) Less(i, j int) bool {
	return strings.Index(s.values[i], s.strValue) < strings.Index(s.values[j], s.strValue)
}

// ForParameter returns the documentation associated with the given parameter name, if any.
func (d SRGDocumentation) ForParameter(paramName string) (SRGDocumentation, bool) {
	strValue, hasValue := documentationForParameter(d.String(), paramName)
	if !hasValue {
		return SRGDocumentation{}, false
	}

	return SRGDocumentation{
		parentComment: d.parentComment,
		isDocComment:  d.isDocComment,
		commentValue:  strValue,
	}, true
}

func documentationForParameter(commentValue string, paramName string) (string, bool) {
	// Find all sentences in the documentation.
	sentences := strings.Split(commentValue, ".")

	// Find any sentences with a ticked reference to the parameter.
	tickedParam := fmt.Sprintf("`%s`", paramName)
	tickedSentences := make([]string, 0, len(sentences))
	for _, sentence := range sentences {
		if strings.Contains(sentence, tickedParam) {
			tickedSentences = append(tickedSentences, sentence)
		}
	}

	if len(tickedSentences) == 0 {
		return "", false
	}

	sort.Sort(byLocation{tickedSentences, tickedParam})
	paramComment := strings.TrimSpace(tickedSentences[0])
	return paramComment, true
}
