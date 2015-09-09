// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Some code based on the Golang VCS portion of the cmd package:
// https://golang.org/src/cmd/go/vcs.go

// vcs package defines helpers functions and interfaces for working with Version Control Systems
// such as git, including discovery of VCS information based on the Golang VCS discovery protocol.
package vcs

// VCSKind identifies the supported kinds of VCS.
type VCSKind int

const (
	VCSKindUnknown VCSKind = iota // an unknown kind of VCS
	VCSKindGit                    // Git
)

// VCSUrlInformation holds information about a VCS source URL.
type VCSUrlInformation struct {
	UrlPrefix    string  // The prefix matching the source URL.
	Kind         VCSKind // The kind of VCS for the source URL.
	DownloadPath string  // The VCS-specific download path.
}

// VCSCache is an interface for a cache that can read and write information about a VCS
// url.
type VCSCache interface {
	// Get returns the cached VCSInformation for the specified URL. If none, returns false
	// an empty struct.
	Get(url string) (VCSUrlInformation, bool)

	// Write saves the given VCSUrlInformation for the given source URL to the cache.
	Write(url string, information VCSUrlInformation)
}

// GetVCSInformation returns the VCS information for a given VCS URL.
func GetVCSInformation(vcsUrl string, cache VCSCache) (VCSUrlInformation, error) {
	// Check the VCS cache first.
	if cache != nil {
		if cached, ok := cache.Get(vcsUrl); ok {
			return cached, nil
		}
	}

	// Perform discovery.
	information, err := discoverInformationForVCSUrl(vcsUrl)
	if err != nil {
		return information, err
	}

	if cache != nil {
		cache.Write(vcsUrl, information)
	}

	return information, err
}

// vcsHandler represents the defined handler information for a specific kind of VCS.
type vcsHandler struct {
	// The kind of the VCS being handled.
	kind VCSKind
}

// vcsById holds a map from string ID for the VCS to its associated vcsHandler struct.
var vcsById = map[string]vcsHandler{
	"git": vcsHandler{
		kind: VCSKindGit,
	},
}
