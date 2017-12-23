// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package integration defines interfaces and helpers for writing integrations with Serulian.
package integration

// IntegrationProviderConstName defines the name of the root-level const that holds the IntegrationsProvider.
const IntegrationProviderConstName = "SerulianIntegrationsProvider"

// IntegrationsProvider defines an interface for providing Integration implementations.
type IntegrationsProvider interface {
	// SerulianIntegrations returns all integrations defined by the provider.
	SerulianIntegrations() []Integration
}

// Integration defines an integration for the Serulian toolkit. Integrations will typically
// match one of the integration interfaces:
//   - LanguageIntegration
type Integration interface{}
