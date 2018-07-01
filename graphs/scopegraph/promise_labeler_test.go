// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/sourceshape"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

type expectedPromiseLabel struct {
	name  string
	label proto.ScopeLabel
}

type promisingLabelTest struct {
	name                  string
	entrypoint            string
	expectedPromiseLabels []expectedPromiseLabel
}

var promisingLabelTests = []promisingLabelTest{
	promisingLabelTest{"empty function no promise test", "emptyfunction",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"EmptyFunction", proto.ScopeLabel_SML_PROMISING_NO},
		},
	},

	promisingLabelTest{"called function no promise test", "calledfunction",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"CalledFunction", proto.ScopeLabel_SML_PROMISING_NO},
			expectedPromiseLabel{"CallingFunction", proto.ScopeLabel_SML_PROMISING_NO},
		},
	},

	promisingLabelTest{"dynamic access no promise test", "dynamicaccessnopromise",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"DynamicAccess", proto.ScopeLabel_SML_PROMISING_NO},
		},
	},

	promisingLabelTest{"dynamic access maybe promise test", "dynamicaccesspromise",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"DynamicAccess", proto.ScopeLabel_SML_PROMISING_MAYBE},
		},
	},

	promisingLabelTest{"dynamic access field promise test", "dynamicaccessfield",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"doSomething", proto.ScopeLabel_SML_PROMISING_NO},
		},
	},

	promisingLabelTest{"unknown function called maybe promise test", "unknownfunc",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"CallsUnknownFunc", proto.ScopeLabel_SML_PROMISING_MAYBE},
		},
	},

	promisingLabelTest{"awaiting promise test", "awaiting",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"Awaiting", proto.ScopeLabel_SML_PROMISING_YES},
			expectedPromiseLabel{"CallsAwaiting", proto.ScopeLabel_SML_PROMISING_YES},
		},
	},

	promisingLabelTest{"awaiting interface promise test", "awaitinginterface",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"AnotherFunction", proto.ScopeLabel_SML_PROMISING_MAYBE},
		},
	},

	promisingLabelTest{"interface no promise test", "interfacenonpromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"AnotherFunction", proto.ScopeLabel_SML_PROMISING_NO},
		},
	},

	promisingLabelTest{"basic sync no promise test", "basicsync",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"DoSomething", proto.ScopeLabel_SML_PROMISING_NO},
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_NO},
		},
	},

	promisingLabelTest{"property not promising test", "propertynotpromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"SomeClass.SomeProp.get", proto.ScopeLabel_SML_PROMISING_NO},
			expectedPromiseLabel{"SomeClass.SomeProp.set", proto.ScopeLabel_SML_PROMISING_NO},

			expectedPromiseLabel{"AnotherFunction", proto.ScopeLabel_SML_PROMISING_NO},
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_NO},
		},
	},

	promisingLabelTest{"property promising test", "propertypromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"SomeClass.SomeProp.get", proto.ScopeLabel_SML_PROMISING_YES},
			expectedPromiseLabel{"SomeClass.SomeProp.set", proto.ScopeLabel_SML_PROMISING_NO},

			expectedPromiseLabel{"AnotherFunction", proto.ScopeLabel_SML_PROMISING_YES},
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_YES},
		},
	},

	promisingLabelTest{"interface property promising test", "interfacepropertypromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"SomeClass.SomeProp.get", proto.ScopeLabel_SML_PROMISING_YES},
			expectedPromiseLabel{"SomeClass.SomeProp.set", proto.ScopeLabel_SML_PROMISING_NO},

			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_MAYBE},
		},
	},

	promisingLabelTest{"loop stream not promising test", "loopstreamnotpromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_NO},
		},
	},

	promisingLabelTest{"loop stream maybe promising test", "loopstreammaybepromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_MAYBE},
		},
	},

	promisingLabelTest{"with statement not promising test", "withnotpromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_NO},
		},
	},

	promisingLabelTest{"with statement promising test", "withpromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_YES},
		},
	},

	promisingLabelTest{"with statement indirect promising test", "withindirectpromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_YES},
		},
	},

	promisingLabelTest{"generator promising test", "generatorpromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"SomeGenerator", proto.ScopeLabel_SML_PROMISING_YES},
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_MAYBE},
		},
	},

	promisingLabelTest{"generator not promising test", "generatornotpromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"SomeGenerator", proto.ScopeLabel_SML_PROMISING_YES},
			expectedPromiseLabel{"AnotherGenerator", proto.ScopeLabel_SML_PROMISING_NO},
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_MAYBE},
		},
	},

	promisingLabelTest{"generator next maybe promising test", "generatornextmaybepromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"SomeGenerator", proto.ScopeLabel_SML_PROMISING_YES},
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_MAYBE},
		},
	},

	promisingLabelTest{"nullable invoke promising test", "nullableinvoke",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"SomeClass.SomeFunction", proto.ScopeLabel_SML_PROMISING_YES},
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_MAYBE},
		},
	},

	promisingLabelTest{"nullable invoke maybe promising test", "nullablemaybeinvoke",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"SomeClass.SomeFunction", proto.ScopeLabel_SML_PROMISING_MAYBE},
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_MAYBE},
		},
	},

	promisingLabelTest{"static function maybe promising test", "staticmaybepromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"SomeClass.DoSomething", proto.ScopeLabel_SML_PROMISING_YES},
			expectedPromiseLabel{"Caller", proto.ScopeLabel_SML_PROMISING_MAYBE},
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_MAYBE},
		},
	},

	promisingLabelTest{"template string not promising test", "templatestringnotpromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_NO},
		},
	},

	promisingLabelTest{"tagged template string promising test", "taggedtemplatestringpromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"SomeTemplateString", proto.ScopeLabel_SML_PROMISING_YES},
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_YES},
		},
	},

	promisingLabelTest{"async var does not cause promising test", "asyncvarnotpromising",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"someVar", proto.ScopeLabel_SML_PROMISING_YES},
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_NO},
		},
	},

	promisingLabelTest{"nested sml children promising test", "nestedsml",
		[]expectedPromiseLabel{
			expectedPromiseLabel{"anotherfunction", proto.ScopeLabel_SML_PROMISING_NO},
			expectedPromiseLabel{"somefunction", proto.ScopeLabel_SML_PROMISING_NO},
			expectedPromiseLabel{"TEST", proto.ScopeLabel_SML_PROMISING_MAYBE},
		},
	},
}

func TestPromisingLabels(t *testing.T) {
	defer compilerutil.DetectGoroutineLeak(t, runtime.NumGoroutine())

	for _, test := range promisingLabelTests {
		if os.Getenv("FILTER") != "" && !strings.Contains(test.name, os.Getenv("FILTER")) {
			continue
		}

		fmt.Printf("Running promise label test: %v\n", test.name)

		entrypointFile := "tests/promising/" + test.entrypoint + ".seru"
		result, _ := ParseAndBuildScopeGraph(entrypointFile, []string{}, packageloader.Library{TESTLIB_PATH, false, "", "testcore"})
		if !assert.True(t, result.Status, "Expected success in scoping on test: %v\n%v\n%v", test.name, result.Errors, result.Warnings) {
			continue
		}

		// Check the promising labels.
		for _, expectedLabel := range test.expectedPromiseLabels {
			pieces := strings.Split(expectedLabel.name, ".")

			member, found := result.Graph.TypeGraph().LookupModuleMember(expectedLabel.name, compilercommon.InputSource(entrypointFile))
			if len(pieces) > 1 {
				parentType, found := result.Graph.TypeGraph().LookupType(pieces[0], compilercommon.InputSource(entrypointFile))
				if !assert.True(t, found, "Missing type %s in test: %v", pieces[0], test.name) {
					continue
				}

				member, found = parentType.GetMember(pieces[1])
				if !assert.True(t, found, "Missing member %s under type %s in test: %v", pieces[1], parentType.Name(), test.name) {
					continue
				}
			} else {
				if !assert.True(t, found, "Missing member %s in test: %v", expectedLabel.name, test.name) {
					continue
				}
			}

			sourceNodeId, _ := member.SourceNodeId()
			sourceNode := result.Graph.SourceGraph().GetNode(sourceNodeId)

			if len(pieces) > 2 {
				if pieces[2] == "get" {
					sourceNode = sourceNode.GetNode(sourceshape.NodePropertyGetter)
				} else if pieces[2] == "set" {
					sourceNode = sourceNode.GetNode(sourceshape.NodePropertySetter)
				}
			}

			sli := result.Graph.layer.
				StartQuery(sourceNodeId).
				In(NodePredicateLabelSource).
				BuildNodeIterator()

			secondaryLabels := make([]proto.ScopeLabel, 0)
			for sli.Next() {
				value, _ := strconv.Atoi(sli.Node().Get(NodePredicateSecondaryLabelValue))
				if value > 0 {
					secondaryLabels = append(secondaryLabels, proto.ScopeLabel(value))
				}
			}

			if !assert.True(t, result.Graph.HasSecondaryLabel(sourceNode, expectedLabel.label), "Expected promising label %v on node %s in test: %v. Found: %v", expectedLabel.label, expectedLabel.name, test.name, secondaryLabels) {
				continue
			}
		}
	}
}
