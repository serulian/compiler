// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import "github.com/serulian/compiler/graphs/scopegraph/proto"

// statementLabelSet defines a set of labels on a statement.
type statementLabelSet struct {
	labelsMap map[proto.ScopeLabel]int
	labels    []proto.ScopeLabel
}

// newLabelSet returns a new statement label set.
func newLabelSet() *statementLabelSet {
	return &statementLabelSet{
		labelsMap: map[proto.ScopeLabel]int{},
		labels:    make([]proto.ScopeLabel, 0),
	}
}

// AppendLabelsOf appends any propagating statement labels found to this set.
func (sls *statementLabelSet) AppendLabelsOf(scope *proto.ScopeInfo) {
	for _, label := range scope.Labels {
		sls.Append(label)
	}
}

// Append appends the given label to the set.
func (sls *statementLabelSet) Append(label proto.ScopeLabel) {
	if _, found := sls.labelsMap[label]; !found {
		sls.labelsMap[label] = len(sls.labels)
		sls.labels = append(sls.labels, label)
	}
}

// HasLabel returns whether the label set contains the specified label.
func (sls *statementLabelSet) HasLabel(expected proto.ScopeLabel) bool {
	_, ok := sls.labelsMap[expected]
	return ok
}

// GetLabels returns all the labels in this set.
func (sls *statementLabelSet) GetLabels() []proto.ScopeLabel {
	return sls.labels
}

// RemoveLabels removes the given labels from the set.
func (sls *statementLabelSet) RemoveLabels(labels ...proto.ScopeLabel) {
	for _, label := range labels {
		index, ok := sls.labelsMap[label]
		if !ok {
			continue
		}

		sls.labels = append(sls.labels[:index], sls.labels[index+1:]...)
		delete(sls.labelsMap, label)
	}
}
