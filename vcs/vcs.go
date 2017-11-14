// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Some code based on the Golang VCS portion of the cmd package:
// https://golang.org/src/cmd/go/vcs.go

// vcs package defines helpers functions and interfaces for working with Version Control Systems
// such as git, including discovery of VCS information based on the Golang VCS discovery protocol.
//
// The discovery protocol is documented loosely here: https://golang.org/cmd/go/#hdr-Remote_import_paths
package vcs

import (
	"fmt"
	"log"
	"os"
	"path"
)

// VCSPackageStatus is the status of the VCS package checked out.
type VCSPackageStatus int

const (
	// DetachedPackage indicates that the package is detatched from a branch and therefore
	// is static.
	DetachedPackage VCSPackageStatus = iota

	// BranchOrHEADPackage indicates that the package is a branch or head package, and will
	// therefore be updated on every call.
	BranchOrHEADPackage

	// LocallyModifiedPackage indicates that the package was modified on the local file system,
	// and therefore cannot be updated.
	LocallyModifiedPackage

	// DevelopmentPackage indicates that the package was found in the VCS development directory
	// and was therefore loaded from that location.
	DevelopmentPackage

	// CachedPackage indicates that the package was returned from cache without further operation.
	// Should only be returned if the always-use-cache options is specified (typically by tooling).
	CachedPackage
)

// VCSCheckoutResult is the result of a VCS checkout, if it succeeds.
type VCSCheckoutResult struct {
	PackageDirectory string
	Warning          string
	Status           VCSPackageStatus
}

// IsVCSRootDirectory returns true if the given local file system path is a VCS root directory.
// Note that this method will return false if the path does not exist locally.
func IsVCSRootDirectory(localPath string) bool {
	_, hasHandler := DetectHandler(localPath)
	return hasHandler
}

// GetVCSCheckoutDirectory returns the path of the directory into which the given VCS path will checked out,
// if PerformVCSCheckout is called.
func GetVCSCheckoutDirectory(vcsPath string, pkgCacheRootPath string, vcsDevelopmentDirectories ...string) (string, error) {
	// Parse the VCS path.
	parsedPath, perr := ParseVCSPath(vcsPath)
	if perr != nil {
		return "", perr
	}

	// Check for a local development cache.
	if parsedPath.isRepoOnlyReference() {
		for _, directoryPath := range vcsDevelopmentDirectories {
			fullCheckDirectory := path.Join(directoryPath, parsedPath.url)
			if _, err := os.Stat(fullCheckDirectory); err == nil {
				return fullCheckDirectory, nil
			}
		}
	}

	fullCacheDirectory := path.Join(pkgCacheRootPath, parsedPath.cacheDirectory())
	return fullCacheDirectory, nil
}

// VCSCacheOption defines the caching options for VCS checkout.
type VCSCacheOption int

const (
	// VCSFollowNormalCacheRules indicates that VCS checkouts will be pulled from cache unless a HEAD
	// reference.
	VCSFollowNormalCacheRules VCSCacheOption = iota

	// VCSAlwaysUseCache indicates that VCS checkouts will always use the cache if available.
	VCSAlwaysUseCache
)

// PerformVCSCheckout performs the checkout and updating of the given VCS path and returns
// the local system directory at which the package was checked out.
//
// pkgCacheRootPath holds the path of the root directory that forms the package cache.
//
// vcsDevelopmentDirectories specifies optional directories to check for branchless and tagless copies
// of the repository first. If found, the copy will be used in lieu of a normal checkout.
func PerformVCSCheckout(vcsPath string, pkgCacheRootPath string, cacheOption VCSCacheOption, vcsDevelopmentDirectories ...string) (VCSCheckoutResult, error) {
	// Parse the VCS path.
	parsedPath, perr := ParseVCSPath(vcsPath)
	if perr != nil {
		return VCSCheckoutResult{}, perr
	}

	var err error
	var warning string

	// Check for a local development cache.
	if parsedPath.isRepoOnlyReference() {
		for _, directoryPath := range vcsDevelopmentDirectories {
			fullCheckDirectory := path.Join(directoryPath, parsedPath.url)
			if _, err := os.Stat(fullCheckDirectory); err == nil {
				// Found a local copy. Returning it.
				log.Printf("Found a local HEAD copy of package %s under development directory %s", parsedPath.url, directoryPath)
				warning = fmt.Sprintf(
					`Package '%s' does not specify a tag, commit or branch and a local copy was found under development directory '%s'. VCS checkout will be skipped and the local copy used instead. To return to normal VCS behavior for this package, remove the --vcs-dev-dir flag or specify the package's tag, commit or branch.`, parsedPath.String(), directoryPath)
				return VCSCheckoutResult{fullCheckDirectory, warning, DevelopmentPackage}, nil
			}
		}
	}

	// Conduct the checkout or update.
	fullCacheDirectory := path.Join(pkgCacheRootPath, parsedPath.cacheDirectory())
	status, err := checkCacheAndPull(parsedPath, fullCacheDirectory, cacheOption)
	if err != nil {
		return VCSCheckoutResult{}, err
	}

	// Warn if the package is HEAD checkout.
	if status == BranchOrHEADPackage {
		warning = fmt.Sprintf("Package '%s' points to HEAD or a branch and will be updated on every build", parsedPath.String())
	}

	// If the parsed path is a subdirectory of the checkout, return a reference to it.
	if parsedPath.subpackage != "" {
		subpackageCacheDirectory := path.Join(fullCacheDirectory, parsedPath.subpackage)
		if _, serr := os.Stat(subpackageCacheDirectory); os.IsNotExist(serr) {
			return VCSCheckoutResult{}, fmt.Errorf("Subpackage '%s' does not exist under VCS package '%s'", parsedPath.subpackage, parsedPath.url)
		}

		return VCSCheckoutResult{subpackageCacheDirectory, warning, status}, nil
	}

	return VCSCheckoutResult{fullCacheDirectory, warning, status}, nil
}

type tagGetter func() ([]string, error)

// InspectInfo holds all the data returned from a call to PerformVCSCheckoutAndInspect.
type InspectInfo struct {
	Engine    string
	CommitSHA string
	GetTags   tagGetter
}

// PerformVCSCheckoutAndInspect performs the checkout and updating of the given VCS path and returns
// the commit SHA of the package, as well as its tags.
//
// pkgCacheRootPath holds the path of the root directory that forms the package cache.
//
// vcsDevelopmentDirectories specifies optional directories to check for branchless and tagless copies
// of the repository first. If found, the copy will be used in lieu of a normal checkout.
func PerformVCSCheckoutAndInspect(vcsPath string, pkgCacheRootPath string, cacheOption VCSCacheOption, vcsDevelopmentDirectories ...string) (InspectInfo, string, error) {
	result, err := PerformVCSCheckout(vcsPath, pkgCacheRootPath, cacheOption, vcsDevelopmentDirectories...)
	if err != nil {
		return InspectInfo{}, "", err
	}

	handler, _ := DetectHandler(result.PackageDirectory)
	sha, err := handler.Inspect(result.PackageDirectory)
	if err != nil {
		return InspectInfo{}, "", err
	}

	tagGetter := func() ([]string, error) {
		return handler.ListTags(result.PackageDirectory)
	}

	return InspectInfo{handler.Kind(), sha, tagGetter}, result.Warning, nil
}

// checkCacheAndPull conducts the cache check and necessary pulls.
func checkCacheAndPull(parsedPath vcsPackagePath, fullCacheDirectory string, cacheOption VCSCacheOption) (VCSPackageStatus, error) {
	// Check the package cache for the path.
	log.Printf("Checking cache directory %s", fullCacheDirectory)
	if _, err := os.Stat(fullCacheDirectory); os.IsNotExist(err) {
		// Do a full checkout.
		return performFullCheckout(parsedPath, fullCacheDirectory)
	}

	log.Printf("Cache directory %s exists", fullCacheDirectory)

	// If caching should always be used, just return the cached package.
	if cacheOption == VCSAlwaysUseCache {
		return CachedPackage, nil
	}

	return performPossibleUpdateCheckout(parsedPath, fullCacheDirectory)
}

// performFullCheckout performs a full VCS checkout of the given package path.
func performFullCheckout(path vcsPackagePath, fullCacheDirectory string) (VCSPackageStatus, error) {
	// Lookup the VCS discovery information.
	discovery, err := DiscoverVCSInformation(path.url)
	if err != nil {
		return DetachedPackage, err
	}

	// Perform a full checkout.
	handler, ok := GetHandlerByKind(discovery.Kind)
	if !ok {
		panic("Unknown VCS handler")
	}

	err = handler.Checkout(path, discovery.DownloadPath, fullCacheDirectory)
	if err != nil {
		return DetachedPackage, err
	}

	// Check if it currently points to a detached state.
	isDetached, err := handler.IsDetached(fullCacheDirectory)
	if err != nil {
		return DetachedPackage, err
	}

	if isDetached {
		return DetachedPackage, nil
	}

	return BranchOrHEADPackage, nil
}

// performPossibleUpdateCheckout performs a VCS update of the given package path, if necessary.
func performPossibleUpdateCheckout(path vcsPackagePath, fullCacheDirectory string) (VCSPackageStatus, error) {
	// Detect the kind of VCS based on the checkout.
	handler, ok := DetectHandler(fullCacheDirectory)
	if !ok {
		return DetachedPackage, fmt.Errorf("Could not detect VCS for directory: %s", fullCacheDirectory)
	}

	// Check for any local file changes.
	if handler.HasLocalChanges(fullCacheDirectory) {
		log.Printf("Package %s is locally modified", fullCacheDirectory)
		return LocallyModifiedPackage, nil
	}

	// Check for a detached package. If found, nothing more to do.
	isDetached, err := handler.IsDetached(fullCacheDirectory)
	if isDetached || err != nil {
		return DetachedPackage, err
	}

	return BranchOrHEADPackage, handler.Update(fullCacheDirectory)
}
