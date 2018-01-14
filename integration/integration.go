// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package integration defines interfaces and helpers for writing integrations with Serulian.
package integration

import (
	"path"
)

// IntegrationConstName defines the name of the root-level const that holds the Integration.
const IntegrationConstName = "SerulianIntegration"

// Integration defines an interface for returing all integration implementations.
type Integration interface {
	// SerulianIntegrations returns all integrations defined by the provider.
	IntegrationImplementations() []Implementation
}

// Implementation defines an integration for the Serulian toolkit. Implementation will typically
// match one of the integration interfaces:
//   - LanguageIntegration
type Implementation interface{}

type IntegrationInformation struct {
	path        string
	integration Integration
}

func (ii IntegrationInformation) Integration() Integration {
	return ii.integration
}

func (ii IntegrationInformation) ID() string {
	return path.Base(ii.path)
}

func (ii IntegrationInformation) matchesID(integrationID string) bool {
	return ii.ID() == integrationID || ii.ID() == integrationID+integrationSuffix
}

func (ii IntegrationInformation) manifest() (integrationManifest, error) {
	return loadManifest(ii.manifestPath())
}

func (ii IntegrationInformation) manifestPath() string {
	return ii.path + ".json"
}

func (ii IntegrationInformation) binaryPath() string {
	return ii.path
}
