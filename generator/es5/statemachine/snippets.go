// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"fmt"

	"github.com/serulian/compiler/generator/es5/shared"
)

var _ = fmt.Sprint

type snippets struct {
	functionTraits shared.StateFunctionTraits
	templater      *shared.Templater
	resources      *ResourceStack
}

func (s snippets) GeneratorDone() string {
	template := `
		$done();
		return;
	`

	return s.templater.Execute("generatordone", template, nil)
}

func (s snippets) Resolve(value string) string {
	var template = ""
	switch s.functionTraits.Type() {
	case shared.StateFunctionNormalSync:
		template = `
			return {{ . }};
		`

		if s.functionTraits.ManagesResources() {
			template = `
				{{ if . }}
				var $pat = {{ . }};
				{{ else }}
				var $pat = undefined;
				{{ end }}
				$resources.popall();
				return $pat;
			`
		}

	case shared.StateFunctionNormalAsync:
		template = `
			$resolve({{ . }});
			return;
		`

	default:
		panic("Unknown function kind in resolve")
	}

	return s.templater.Execute("resolve", template, value)
}

func (s snippets) Reject(value string) string {
	var template = ""
	if s.functionTraits.IsAsynchronous() {
		template = `
			$reject({{ . }});
			return;
		`
	} else {
		template = `
			throw {{ . }};
		`

		if s.functionTraits.ManagesResources() {
			template = `
				var $pat = {{ . }};
				$resources.popall();
				throw $pat;
			`
		}
	}

	return s.templater.Execute("reject", template, value)
}

func (s snippets) Continue(isUnderPromise bool) string {
	var template = ""
	switch s.functionTraits.Type() {
	case shared.StateFunctionNormalSync:
		template = `
			continue syncloop;
		`

	case shared.StateFunctionNormalAsync:
		if isUnderPromise {
			template = `
				$continue($resolve, $reject);
				return;
			`
		} else {
			template = `				
				continue localasyncloop;
			`
		}

	case shared.StateFunctionSyncOrAsyncGenerator:
		template = `
			$continue($yield, $yieldin, $reject, $done);
			return;
		`

	default:
		panic("Unknown function kind in continue")
	}

	return s.templater.Execute("continue", template, nil)
}

func (s snippets) setState(id stateId, followingAction string) string {
	data := struct {
		TargetID        int
		FollowingAction string
		ResourceNames   []string
		IsAsync         bool
	}{int(id), followingAction, s.resources.OutOfScope(id), s.functionTraits.IsAsynchronous()}

	template := `
		$current = {{ .TargetID }};
		{{ .FollowingAction }}
	`

	if len(data.ResourceNames) > 0 {
		template = `
		  $resources.popr(
			{{ range $idx, $name := .ResourceNames }}
				{{ if $idx }},{{ end }}
				'{{ $name }}'
			{{ end }}
		  ){{ if .IsAsync }}.then(function() {
		  	$current = {{ .TargetID }};
		  	{{ .FollowingAction }}
		  }){{ end }};

		  {{ if not .IsAsync }}
			$current = {{ .TargetID }};
			{{ .FollowingAction }}
		  {{ end }}
		`
	}

	return s.templater.Execute("setstate", template, data)
}

func (s snippets) SetStateAndBreak(state *state) string {
	return s.setState(state.ID, "return")
}

func (s snippets) SetStateAndContinue(state *state, isUnderPromise bool) string {
	return s.setState(state.ID, s.Continue(isUnderPromise))
}
