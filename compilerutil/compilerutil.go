// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// compilerutil package defines utility methods.
package compilerutil

import (
	"fmt"
	"os"

	"github.com/nu7hatch/gouuid"
)

type checkFn func() bool

// DCHECK executes the checker function only if the DEBUG environment variable is set.
// If the function returns false, the compiler will panic with the formatted message.
func DCHECK(checker checkFn, failMessage string, args ...interface{}) {
	if os.Getenv("DEBUG") == "" {
		return
	}

	if !checker() {
		panic(fmt.Sprintf(failMessage, args...))
	}
}

// IsId returns whether the given string is a possible ID as returned by NewUniqueId.
func IsId(possibleId string) bool {
	_, err := uuid.ParseHex(possibleId)
	return err == nil
}

// NewUniqueId returns a new unique ID.
func NewUniqueId() string {
	u4, err := uuid.NewV4()
	if err != nil {
		panic(err)
		return ""
	}

	return u4.String()
}
