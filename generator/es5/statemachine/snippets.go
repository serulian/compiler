// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"github.com/serulian/compiler/generator/es5/shared"
)

type snippets struct {
	templater *shared.Templater
}

func (s snippets) Resolve(value string) string {
	template := `
		$resolve({{ . }});
		return;
	`

	return s.templater.Execute("resolve", template, value)
}

func (s snippets) Reject(value string) string {
	template := `
		$reject({{ . }});
		return;
	`

	return s.templater.Execute("reject", template, value)
}

func (s snippets) Continue() string {
	template := `
		$continue($resolve, $reject);
		return;
	`

	return s.templater.Execute("continue", template, nil)
}

func (s snippets) SetState(id stateId) string {
	template := `
		$current = {{ . }};
	`

	return s.templater.Execute("setstate", template, int(id))
}
