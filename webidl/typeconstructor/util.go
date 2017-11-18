// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typeconstructor

import (
	"path"
	"strings"
)

// determineSharedRootDirectory returns the shared root directory for all the path strings
// given. If there is none, returns an empty string.
func determineSharedRootDirectory(paths []string) string {
	if len(paths) == 0 {
		return ""
	}

	var commonPath = path.Dir(paths[0])
	for _, current := range paths[1:] {
		cleaned := path.Dir(current)
		commonPath = commonPieces(cleaned, commonPath, "/")
		if len(commonPath) == 0 {
			break
		}
	}

	if len(commonPath) == 0 {
		return ""
	}

	return commonPath + "/"
}

// commonPieces returns the common prefix pieces of the given strings, broken into
// pieces using the sep separator.
func commonPieces(first string, second string, sep string) string {
	firstPieces := strings.Split(first, sep)
	secondPieces := strings.Split(second, sep)

	var length = len(firstPieces)
	if len(secondPieces) < len(firstPieces) {
		length = len(secondPieces)
	}

	var collected = make([]string, 0, length)
	for i := 0; i < length; i++ {
		if firstPieces[i] != secondPieces[i] {
			break
		}

		collected = append(collected, firstPieces[i])
	}

	return strings.Join(collected, sep)
}
