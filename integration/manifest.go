// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

const manifestFileName = "serulian-integration.json"
const manifestVersion = "v1"

// integrationManifest defines the structure of the integration metadata found in the
// manifest file at the root of a provider's repository.
type integrationManifest struct {
	ManifestVersion string `json:"manifestVersion"`

	Name    string `json:"name"`
	Title   string `json:"title"`
	Version string `json:"version"`
	RepoURL string `json:"repoURL"`

	ProjectURL  string `json:"projectURL,omitempty"`
	Description string `json:"description,omitempty"`

	Download   downloadManifest        `json:"download"`
	VersionMap map[string]versionEntry `json:"versionMap"`
}

func (im integrationManifest) JSONString() ([]byte, error) {
	return json.Marshal(im)
}

// getIntegrationVersionForToolkit returns the latest version matching the specified toolkit version.
func (im integrationManifest) getIntegrationVersionForToolkit(toolkitVersion string, versions []string) (string, bool) {
	supportedVersions := make([]string, 0)
	for _, version := range versions {
		entry, ok := im.VersionMap[version]
		if !ok {
			continue
		}

		if entry.ToolkitVersion == toolkitVersion {
			supportedVersions = append(supportedVersions, version)
		}
	}

	if len(supportedVersions) == 0 {
		return "", false
	}

	// TODO: sort this using semver.
	sort.Strings(supportedVersions)
	return supportedVersions[len(supportedVersions)-1], true
}

// getIntegrationDownloadURL returns the download URL for the integration, matching the given integration version, arch and OS.
func (im integrationManifest) getIntegrationDownloadURL(version string, arch string, os string) (string, error) {
	template := im.Download.Template
	if template == "" {
		return "", fmt.Errorf("Integration download template is invalid")
	}

	template = strings.Replace(template, "{version}", version, -1)
	template = strings.Replace(template, "{arch}", arch, -1)
	template = strings.Replace(template, "{os}", os, -1)
	return template, nil
}

// downloadManifest defines metadata on how to download a specific version of the provider.
type downloadManifest struct {
	// template defines the template to use to find binaries for this provider. The template
	// replaces the following variables with their values:
	// `{version}`, `{arch}`, `{os}`
	Template string `json:"template"`
}

// versionEntry defines metadata on how to map a version of the integration to its metadata.
type versionEntry struct {
	ToolkitVersion string `json:"toolkitVersion"`
}

// loadManifest loads the integration manifest found at the given file path.
func loadManifest(filepath string) (integrationManifest, error) {
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		return integrationManifest{}, err
	}

	var manifest integrationManifest
	err = json.Unmarshal(raw, &manifest)
	if err != nil {
		return integrationManifest{}, err
	}

	if manifest.ManifestVersion != manifestVersion {
		return integrationManifest{}, fmt.Errorf("Could not load manifest with version `%s`: The integration is probably out of date", manifest.ManifestVersion)
	}

	return manifest, nil
}
