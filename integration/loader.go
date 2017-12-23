// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/phayes/permbits"
)

// integrationSubDirectory its the subdirectory in which to search for integrations.
const integrationSubDirectory = ".int"

// LoadIntegrationProviders loads all the integration providers found for the current toolkit.
func LoadIntegrationProviders() ([]IntegrationsProvider, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return loadIntegrationProvidersUnderPath(path.Join(dir, integrationSubDirectory))
}

func loadIntegrationProvidersUnderPath(providerDirPath string) ([]IntegrationsProvider, error) {
	_, err := os.Stat(providerDirPath)
	if os.IsNotExist(err) {
		return []IntegrationsProvider{}, nil
	}

	if err != nil {
		return []IntegrationsProvider{}, err
	}

	// Iterate the directory, finding all binaries and trying to load the integrations found within.
	files, err := ioutil.ReadDir(providerDirPath)
	if err != nil {
		return []IntegrationsProvider{}, err
	}

	if len(files) == 0 {
		return []IntegrationsProvider{}, nil
	}

	providers := make([]IntegrationsProvider, 0, len(files))
	for _, f := range files {
		if f.Mode().IsRegular() && !strings.Contains(f.Name(), ".") {
			fullPath := path.Join(providerDirPath, f.Name())
			permissions, err := permbits.Stat(fullPath)
			if err != nil {
				return []IntegrationsProvider{}, err
			}

			if permissions.UserExecute() || permissions.GroupExecute() || permissions.OtherExecute() {
				// Found a binary. Attempt to load the provider from it.
				p, err := plugin.Open(fullPath)
				if err != nil {
					return []IntegrationsProvider{}, err
				}

				providerSymbol, err := p.Lookup(IntegrationProviderConstName)
				if err != nil {
					return []IntegrationsProvider{}, err
				}

				provider, castOk := providerSymbol.(IntegrationsProvider)
				if !castOk {
					return []IntegrationsProvider{}, fmt.Errorf("Could find integration provider in integration `%s`", f.Name())
				}

				providers = append(providers, provider)
			}
		}
	}

	return providers, nil
}
