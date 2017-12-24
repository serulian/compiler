// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/phayes/permbits"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/vcs"
	"github.com/serulian/compiler/version"

	"github.com/olekukonko/tablewriter"
	input "github.com/tcnksm/go-input"
	pb "gopkg.in/cheggaaa/pb.v2"
)

// ListIntegrations lists all installed integrations on the command line.
func ListIntegrations() bool {
	integrations, err := LoadIntegrations()
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not list installed integrations: %v", err)
		return false
	}

	fmt.Println()
	if len(integrations) == 0 {
		fmt.Println("No Serulian integrations installed")
		return true
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Version", "Capabilities"})
	table.SetCenterSeparator("|")
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})

	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor})

	for _, integration := range integrations {
		manifest, err := integration.manifest()
		if err != nil {
			table.Append([]string{integration.ID(), fmt.Sprintf("Error: %v", err), "", ""})
			continue
		}

		capabilities := strings.Join(getCapabilitiesDescription(integration), "\n")
		table.Append([]string{integration.ID(), manifest.Title, manifest.Version, capabilities})
	}

	table.Render()
	fmt.Println()
	return true
}

// DescribeIntegration describes a particular integration.
func DescribeIntegration(integrationID string) bool {
	integration, err := lookupIntegrationByID(integrationID)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not lookup integeration: %v", err)
		return false
	}

	manifest, err := integration.manifest()
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not read integeration manifest: %v", err)
		return false
	}

	fmt.Println()

	header := fmt.Sprintf("%s v%s", manifest.Title, manifest.Version)
	fmt.Printf("%s\n%s\n", header, strings.Repeat("-", len(header)))

	if len(manifest.Description) > 0 {
		fmt.Printf("%s\n\n", manifest.Description)
	}

	if len(manifest.ProjectURL) > 0 {
		fmt.Printf("Website: %s\n", manifest.ProjectURL)
	}

	fmt.Println()
	fmt.Printf("Capabilities:\n")
	for _, capability := range getCapabilitiesDescription(integration) {
		fmt.Printf("  - %s\n", capability)
	}

	fmt.Println()
	return true
}

// InstallIntegration installs a particular integration.
func InstallIntegration(repoURL string, debug bool) bool {
	// Disable logging unless the debug flag is on.
	if !debug {
		log.SetOutput(ioutil.Discard)
	}

	manifest, versions, ok := checkoutIntegrationManifestAndVersions(repoURL)
	if !ok {
		return false
	}

	// Find the latest released version for the toolkit version.
	integrationVersion, ok := manifest.getIntegrationVersionForToolkit(version.Version, versions)
	if !ok {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "No supported version of the integration was found for this version of the toolkit")
		return false
	}

	return installIntegrationVersion(manifest, integrationVersion)
}

func checkoutIntegrationManifestAndVersions(repoURL string) (integrationManifest, []string, bool) {
	parsedPath, err := vcs.ParseVCSPath(repoURL)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not find a Serulian integration at the specified URL: %v", err)
		return integrationManifest{}, []string{}, false
	}

	compilerutil.LogToConsole(compilerutil.InfoLogLevel, nil, "Looking up integration under repository `%s`", repoURL)
	info, err := vcs.DiscoverVCSInformation(repoURL)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not find a Serulian integration at the specified URL: %v", err)
		return integrationManifest{}, []string{}, false
	}

	// Perform a full checkout into a temp directory.
	handler, ok := vcs.GetHandlerByKind(info.Kind)
	if !ok {
		panic("Unknown VCS handler")
	}

	dir, err := ioutil.TempDir("", "integration")
	if err != nil {
		log.Fatal(err)
	}

	compilerutil.LogToConsole(compilerutil.InfoLogLevel, nil, "Checking out integration under repository `%s`", repoURL)
	err = handler.Checkout(parsedPath, info.DownloadPath, dir)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not checkout the Serulian integration at the specified URL: %v", err)
		return integrationManifest{}, []string{}, false
	}

	// Find the manifest.
	manifest, err := loadManifest(path.Join(dir, manifestFileName))
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not read the Serulian integration manifest: %v", err)
		return integrationManifest{}, []string{}, false
	}

	// Load the released versions.
	versions, err := handler.ListTags(dir)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not read the Serulian integration versions: %v", err)
		return integrationManifest{}, []string{}, false
	}

	return manifest, versions, true
}

func installIntegrationVersion(manifest integrationManifest, integrationVersion string) bool {
	// Get the download URL.
	downloadURL, err := manifest.getIntegrationDownloadURL(integrationVersion, runtime.GOARCH, runtime.GOOS)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "%v", err)
		return false
	}

	// Make sure the URL is secure.
	url, err := url.ParseRequestURI(downloadURL)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not parse download URL: %v", err)
		return false
	}

	if url.Scheme != "https" {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Cannot download from insecure URL: %s", downloadURL)
		return false
	}

	// Download the integration binary.
	integrationPath := manifest.Name + integrationSuffix
	var bar *pb.ProgressBar
	err = compilerutil.DownloadURLToFileWithProgress(url, integrationPath, func(bytesDownloaded int64, totalBytes int64, startCall bool) {
		if startCall {
			bar = pb.New64(totalBytes)
			bar.Set(pb.Bytes, true)
			bar.SetRefreshRate(time.Second)
			bar.Start()
		}

		bar.SetCurrent(bytesDownloaded)
	})

	if bar != nil {
		bar.Finish()
	}

	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not download integration: %v", err)
		return false
	}

	// Mark the file as being executable.
	permissions, err := permbits.Stat(integrationPath)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not download integration: %v", err)
		return false
	}

	permissions.SetGroupExecute(true)
	permissions.SetOtherExecute(true)

	// Verify the integration.
	integrationInfo, err := loadIntegrationAtPath(integrationPath)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not verify integration: %v", err)
		return false
	}

	// Write out the manifest.
	manifestJSON, err := manifest.JSONString()
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not write integration manifest: %v", err)
		return false
	}

	err = ioutil.WriteFile(integrationInfo.manifestPath(), manifestJSON, 0755)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not write integration manifest: %v", err)
		return false
	}

	return true
}

// UpgradeIntegration upgrades a particular integration by downloading a new version.
func UpgradeIntegration(integrationID string, skipPrompt bool, debug bool) bool {
	// Disable logging unless the debug flag is on.
	if !debug {
		log.SetOutput(ioutil.Discard)
	}

	// Lookup the current integration and its manifest, so we can find the integration repo URL.
	currentManifest, err := loadManifest(integrationID + ".json")
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not load integration manifest: %v", err)
		return false
	}

	// Load the newest copy of the integration off of the remote repository, and find the latest version matching
	// this toolkit version.
	manifest, versions, ok := checkoutIntegrationManifestAndVersions(currentManifest.RepoURL)
	if !ok {
		return false
	}

	// Find the latest released version for the toolkit version.
	integrationVersion, ok := manifest.getIntegrationVersionForToolkit(version.Version, versions)
	if !ok {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "No supported version of the integration was found for this version of the toolkit")
		return false
	}

	if !skipPrompt {
		ui := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}

		query := fmt.Sprintf("Upgrade integration `%s` to version `%s`?\n", integrationID, integrationVersion)
		delete, err := ui.Ask(query, &input.Options{
			Default:  "no",
			Required: true,
			Loop:     true,
		})

		if err != nil {
			log.Fatal(err)
		}

		if delete != "yes" && delete != "y" {
			fmt.Println("Upgrade canceled")
			return false
		}
	}

	return installIntegrationVersion(manifest, integrationVersion)
}

// UninstallIntegration uninstalls a particular integration by deleting its binary.
func UninstallIntegration(integrationID string, skipPrompt bool) bool {
	integration, err := lookupIntegrationByID(integrationID)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Could not lookup integration: %v", err)
		return false
	}

	capabilities := strings.Join(getCapabilitiesDescription(integration), ", ")

	if !skipPrompt {
		ui := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}

		query := fmt.Sprintf("Uninstall integration `%s`?\n\nSerulian will no longer have the following capabilities: %s\n", integrationID, capabilities)
		delete, err := ui.Ask(query, &input.Options{
			Default:  "no",
			Required: true,
			Loop:     true,
		})

		if err != nil {
			log.Fatal(err)
		}

		if delete != "yes" && delete != "y" {
			fmt.Println("Uninstall canceled")
			return false
		}
	}

	os.Remove(integration.binaryPath())
	os.Remove(integration.manifestPath())

	fmt.Printf("Uninstalled integration `%s`\n", integrationID)
	return true
}

func lookupIntegrationByID(integrationID string) (IntegrationInformation, error) {
	integrations, err := LoadIntegrations()
	if err != nil {
		return IntegrationInformation{}, err
	}

	for _, integration := range integrations {
		if integration.matchesID(integrationID) {
			return integration, nil
		}
	}

	return IntegrationInformation{}, fmt.Errorf("Could not find integration with ID `%s`", integrationID)
}

func getCapabilitiesDescription(integrationInfo IntegrationInformation) []string {
	var capabilityDescriptions = []string{}
	for _, implementation := range integrationInfo.integration.IntegrationImplementations() {
		if langIntegration, isLangIntegration := implementation.(LanguageIntegration); isLangIntegration {
			capabilityDescriptions = append(capabilityDescriptions, fmt.Sprintf(".%s source files", langIntegration.SourceHandler().PackageFileExtension()))
		}
	}
	return capabilityDescriptions
}
