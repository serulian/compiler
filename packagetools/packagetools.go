// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package packagetools implements tools for working on Serulian packages.
package packagetools

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/serulian/compiler/builder"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/graphs/typegraph/diff"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/vcs"
)

// Revise tags the current package with a semantic version computed by performing a diff of the package's current
// contents against the version found in the same repository with (optional) comparison
// version. If the comparison version is not specified, then the latest semver
// of the package is used (if any).
func Revise(packagePath string, comparisonVersion string, revisionNote string, verbose bool, debug bool, vcsDevelopmentDirectories ...string) {
	// Perform the diff of the package, outputting if necessary. This will also verify the package is valid.
	packageDiff, upstreamVersion, ok := getAndOutputDiff(packagePath, comparisonVersion, verbose, debug, vcsDevelopmentDirectories)
	if !ok {
		return
	}

	fmt.Println()

	// Ensure there are no local uncommitted changes in the package.
	handler, _ := vcs.DetectHandler(packagePath)
	if handler.HasLocalChanges(packagePath, packageloader.SerulianPackageDirectory) {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil,
			"Package `%s` contains uncommited changes", packagePath)
		return
	}

	// Determine how to revise the semantic version.
	var option = compilerutil.NextPatchVersion

	switch {
	case packageDiff.HasBreakingChange():
		option = compilerutil.NextMajorVersion

	case packageDiff.Kind != diff.Same:
		option = compilerutil.NextMinorVersion
	}

	newVersion, err := compilerutil.NextSemanticVersion(upstreamVersion, option)
	if err != nil {
		log.Fatal(err)
	}

	// Tag the latest commit with the new semantic version.
	if revisionNote == "" {
		revisionNote = "Automatic revision to version " + newVersion
	}

	terr := handler.Tag(packagePath, newVersion, revisionNote)
	if terr != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil,
			"Could not tag package with version `%s`: %v", newVersion, err)
		return
	}

	compilerutil.LogToConsole(compilerutil.SuccessLogLevel, nil,
		"Package `%s` tagged with revised version `%s`", packagePath, newVersion)
}

// OutputDiff outputs the diff between the Serulian package found at the given
// path and the version found in the same repository with (optional) comparison
// version. If the comparison version is not specified, then the latest semver
// of the package is used (if any).
func OutputDiff(packagePath string, comparisonVersion string, verbose bool, debug bool, vcsDevelopmentDirectories ...string) {
	getAndOutputDiff(packagePath, comparisonVersion, verbose, debug, vcsDevelopmentDirectories)
}

func getAndOutputDiff(packagePath string, comparisonVersion string, verbose bool, debug bool, vcsDevelopmentDirectories []string) (diff.PackageDiff, string, bool) {
	// Disable logging unless the debug flag is on.
	if !debug {
		log.SetOutput(ioutil.Discard)
	}

	// Ensure the package is VCS-backed.
	handler, hasHandler := vcs.DetectHandler(packagePath)
	if !hasHandler {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil,
			"Path `%s` is not a valid Serulian package. Please make sure to point to the root of the Serulian package.", packagePath)
		return diff.PackageDiff{}, "", false
	}

	// List the tags and find the latest release or specified comparison version.
	compilerutil.LogToConsole(compilerutil.InfoLogLevel, nil, "Listing versions for package...")

	tags, err := handler.ListTags(packagePath)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil,
			"Could not list tags for package `%s`: %v", packagePath, err)
		return diff.PackageDiff{}, "", false
	}

	if comparisonVersion == "" {
		foundVersion, found := compilerutil.LatestSemanticVersion(tags)
		if !found {
			compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil,
				"Could not find a full semantic version for package `%s`. Package versions: %v", packagePath, tags)
			return diff.PackageDiff{}, "", false
		}

		comparisonVersion = foundVersion
	} else {
		var found = false
		for _, version := range tags {
			if version == comparisonVersion {
				found = true
				break
			}
		}
		if !found {
			compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil,
				"Version `%s` not found for package `%s`. Package versions: %v", comparisonVersion, packagePath, tags)
			return diff.PackageDiff{}, "", false
		}
	}

	// Retrieve the path information for the package.
	pathInfo, err := handler.GetPackagePath(packagePath)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil,
			"Could not get path information for package `%s`: %v", packagePath, err)
		return diff.PackageDiff{}, "", false
	}

	// Checkout the comparison version into a temporary directory.
	dir, err := ioutil.TempDir("", "package-diff")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir)

	compilerutil.LogToConsole(compilerutil.InfoLogLevel, nil, "Cloning version `%s`...", comparisonVersion)
	cerr := handler.Checkout(pathInfo.WithTag(comparisonVersion), pathInfo.URL(), dir)
	if cerr != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil,
			"Could not clone version `%s` for package `%s`: %v", comparisonVersion, packagePath, cerr)
		return diff.PackageDiff{}, "", false
	}

	// Scope and build both versions.
	compilerutil.LogToConsole(compilerutil.InfoLogLevel, nil, "Analyzing version `%s`...", comparisonVersion)
	originalGraph, ok := getTypeGraph(dir, vcsDevelopmentDirectories)
	if !ok {
		return diff.PackageDiff{}, "", false
	}

	compilerutil.LogToConsole(compilerutil.InfoLogLevel, nil, "Analyzing HEAD...")
	updatedGraph, ok := getTypeGraph(packagePath, vcsDevelopmentDirectories)
	if !ok {
		return diff.PackageDiff{}, "", false
	}

	original := diff.TypeGraphInformation{Graph: originalGraph, PackageRootPath: dir}
	updated := diff.TypeGraphInformation{Graph: updatedGraph, PackageRootPath: packagePath}
	filter := func(module typegraph.TGModule) bool {
		// Filter out any VCS loaded packages.
		return !strings.Contains(module.PackagePath(), packageloader.SerulianPackageDirectory)
	}

	compilerutil.LogToConsole(compilerutil.InfoLogLevel, nil, "Computing Diff...")
	computed := diff.ComputeDiff(original, updated, filter)

	fmt.Println()

	for _, packageDiff := range computed.Packages {
		compilerutil.MessageColor.Print(comparisonVersion)
		compilerutil.MessageColor.Print(" → ")
		compilerutil.MessageColor.Print("HEAD")

		if packageDiff.Kind == diff.Same {
			compilerutil.MessageColor.Print(": No changes found\n")
			return packageDiff, comparisonVersion, true
		}

		if packageDiff.HasBreakingChange() {
			compilerutil.MessageColor.Print(": ")
			compilerutil.WarningColor.Print("Breaking changes found\n")
		} else {
			compilerutil.MessageColor.Print(": ")
			compilerutil.SuccessColor.Print("Compatible changes found\n")
		}

		if verbose {
			outputDetailedDiff(packageDiff)
		}

		// We only care about a single package.
		return packageDiff, comparisonVersion, true
	}

	return diff.PackageDiff{}, "", false
}

func outputDetailedDiff(packageDiff diff.PackageDiff) {
	compilerutil.BoldWhiteColor.Print("  Change summary:\n")
	for _, change := range packageDiff.ChangeReason.Expand() {
		compatibility, description := change.Describe()

		if compatibility == diff.BackwardIncompatible {
			compilerutil.ErrorColor.Print("    [Breaking]   ")
		} else {
			compilerutil.SuccessColor.Print("    [Compatible] ")
		}

		fmt.Println(description)
	}

	fmt.Println()
	compilerutil.BoldWhiteColor.Print("  Detailed changes:\n")

	indentation := "    "
	numModified := 0
	for _, typeDiff := range packageDiff.Types {
		switch typeDiff.Kind {
		case diff.Same:
			continue

		case diff.Added:
			fmt.Print(indentation)
			compilerutil.SuccessColor.Print("[+] ")
			fmt.Printf("Type %s\n", typeDiff.Name)

		case diff.Removed:
			fmt.Print(indentation)
			compilerutil.ErrorColor.Print("[X] ")
			fmt.Printf("Type %s\n", typeDiff.Name)

		case diff.Changed:
			fmt.Print(indentation)
			compilerutil.WarningColor.Print("[!] ")
			fmt.Printf("Type %s\n", typeDiff.Name)
			outputDetailedTypeDiff(typeDiff, 6)
		}

		numModified++
	}

	for _, memberDiff := range packageDiff.Members {
		switch memberDiff.Kind {
		case diff.Same:
			continue

		case diff.Added:
			fmt.Print(indentation)
			compilerutil.SuccessColor.Print("[+] ")
			printMemberPathAndName(memberDiff.Updated)

		case diff.Removed:
			fmt.Print(indentation)
			compilerutil.ErrorColor.Print("[X] ")
			printMemberPathAndName(memberDiff.Original)

		case diff.Changed:
			fmt.Print(indentation)
			compilerutil.WarningColor.Print("[!] ")
			printMemberPathAndName(memberDiff.Original)
			outputDetailedMemberDiff(memberDiff, 8)
		}
		numModified++
	}

	if numModified == 0 {
		fmt.Printf("    No modified members or types found under this package")
	}
}

func printMemberPathAndName(member *typegraph.TGMember) {
	compilerutil.FaintWhiteColor.Print(member.Title())
	fmt.Print(" ")
	fmt.Print(member.Name())
	fmt.Println()
}

func outputDetailedTypeDiff(typeDiff diff.TypeDiff, indentCount int) {
	indentation := strings.Repeat(" ", indentCount)

	for _, change := range typeDiff.ChangeReason.Expand() {
		compatibility, description := change.Describe()

		if compatibility == diff.BackwardIncompatible {
			compilerutil.LogMessageToConsole(compilerutil.ErrorColor, indentation+"[Breaking]   ", nil, "%s", description)
		} else {
			compilerutil.LogMessageToConsole(compilerutil.SuccessColor, indentation+"[Compatible] ", nil, "%s", description)
		}
	}

	fmt.Println("  ----- Modified Members ------")
	numModified := 0
	for _, memberDiff := range typeDiff.Members {
		switch typeDiff.Kind {
		case diff.Same:
			continue

		case diff.Added:
			fmt.Print(indentation)
			compilerutil.SuccessColor.Print("[+] ")
			printMemberPathAndName(memberDiff.Updated)

		case diff.Removed:
			fmt.Print(indentation)
			compilerutil.ErrorColor.Print("[X] ")
			printMemberPathAndName(memberDiff.Original)

		case diff.Changed:
			fmt.Print(indentation)
			compilerutil.WarningColor.Print("[!] ")
			printMemberPathAndName(memberDiff.Original)
			outputDetailedMemberDiff(memberDiff, indentCount+8)
		}
		numModified++
	}

	if numModified == 0 {
		fmt.Printf(indentation + "  No modified members found under this type")
	}
}

func outputDetailedMemberDiff(memberDiff diff.MemberDiff, indentCount int) {
	indentation := strings.Repeat(" ", indentCount)

	for _, change := range memberDiff.ChangeReason.Expand() {
		compatibility, description := change.Describe()

		if compatibility == diff.BackwardIncompatible {
			compilerutil.LogMessageToConsole(compilerutil.ErrorColor, indentation+"[Breaking]   ", nil, "%s", description)
		} else {
			compilerutil.LogMessageToConsole(compilerutil.SuccessColor, indentation+"[Compatible] ", nil, "%s", description)
		}
	}
}

func getTypeGraph(packagePath string, vcsDevelopmentDirectories []string) (*typegraph.TypeGraph, bool) {
	scopeResult, err := scopegraph.ParseAndBuildScopeGraph(packagePath, vcsDevelopmentDirectories, builder.CORE_LIBRARY)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "%s", err.Error())
		return nil, false
	}

	if !scopeResult.Status {
		builder.OutputErrors(scopeResult.Errors)
		return nil, false
	}

	return scopeResult.Graph.TypeGraph(), true
}
