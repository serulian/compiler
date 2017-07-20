// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"plugin"
	"strings"

	"github.com/phayes/permbits"
)

// PROVIDER_SYMBOL_NAME defines the name of the symbol to be exported by dynamic language integration provider
// plugins.
const PROVIDER_SYMBOL_NAME = "Provider"

// LoadLanguageIntegrationProviders loads all the language integration providers found under the given
// path. *All* files in the directory without an extension will be treated as a potential provider, so
// callers should make sure that the directory is clean otherwise. Note that if the path does not exist,
// the list returned will be *empty*.
func LoadLanguageIntegrationProviders(providerDirPath string) ([]LanguageIntegrationProvider, error) {
	_, err := os.Stat(providerDirPath)
	if os.IsNotExist(err) {
		return []LanguageIntegrationProvider{}, nil
	}

	if err != nil {
		return []LanguageIntegrationProvider{}, err
	}

	// Iterate the directory, finding all binaries and trying to load the integrations found within.
	files, err := ioutil.ReadDir(providerDirPath)
	if err != nil {
		return []LanguageIntegrationProvider{}, err
	}

	if len(files) == 0 {
		return []LanguageIntegrationProvider{}, nil
	}

	providers := make([]LanguageIntegrationProvider, 0, len(files))
	for _, f := range files {
		if f.Mode().IsRegular() && !strings.Contains(f.Name(), ".") {
			fullPath := path.Join(providerDirPath, f.Name())
			permissions, err := permbits.Stat(fullPath)
			if err != nil {
				return []LanguageIntegrationProvider{}, err
			}

			if permissions.UserExecute() || permissions.GroupExecute() || permissions.OtherExecute() {
				// Found a binary. Attempt to load the provider from it.
				p, err := plugin.Open(fullPath)
				if err != nil {
					return []LanguageIntegrationProvider{}, err
				}

				providerSymbol, err := p.Lookup(PROVIDER_SYMBOL_NAME)
				if err != nil {
					return []LanguageIntegrationProvider{}, err
				}

				provider, castOk := providerSymbol.(LanguageIntegrationProvider)
				if !castOk {
					return []LanguageIntegrationProvider{}, fmt.Errorf("Could find language integration provider in plugin `%s`", f.Name())
				}

				providers = append(providers, provider)
			}
		}
	}

	return providers, nil
}
