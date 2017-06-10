// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"fmt"
	"log"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/formatter"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"
	"github.com/serulian/compiler/vcs"
)

const shortSHALength = 7

// AllActions defines the set of all actions supported by Grok.
var AllActions = []string{string(NoAction), string(FreezeImport), string(UnfreezeImport)}

// ExecuteAction executes an action as defined by GetContextActions.
func (gh Handle) ExecuteAction(action Action, params map[string]interface{}, source compilercommon.InputSource) error {
	switch action {
	case NoAction:
		return nil

	case FreezeImport:
		importSource, hasImportSource := params["import-source"]
		if !hasImportSource {
			return fmt.Errorf("Missing import source")
		}

		tag, hasTag := params["tag"]
		commit, hasCommit := params["commit"]
		commitOrTag := formatter.Tag("0.0.0")
		if hasTag {
			commitOrTag = formatter.Tag(tag.(string))
		} else if hasCommit {
			commitOrTag = formatter.Commit(commit.(string))
		} else {
			return fmt.Errorf("Missing commit or tag")
		}

		if !formatter.FreezeAt(string(source), importSource.(string), commitOrTag, gh.groker.vcsDevelopmentDirectories, false) {
			return fmt.Errorf("Could not freeze import")
		}
		return nil

	case UnfreezeImport:
		importSource, hasImportSource := params["import-source"]
		if !hasImportSource {
			return fmt.Errorf("Missing import source")
		}

		if !formatter.UnfreezeAt(string(source), importSource.(string), gh.groker.vcsDevelopmentDirectories, false) {
			return fmt.Errorf("Could not unfreeze import")
		}
		return nil

	default:
		return fmt.Errorf("Unknown action: %v", action)
	}
}

// GetActionsForPosition returns all asynchronous code actions for the given source position. Unlike GetContextActions, the
// returns items must *all* be actions, and should be displayed in a selector menu, rather than inline in the code.
func (gh Handle) GetActionsForPosition(source compilercommon.InputSource, lineNumber int, colPosition int) ([]ContextOrAction, error) {
	// Note: We cannot use PositionFromLineAndColumn here, as it will lookup the position in the *tracked* source, which
	// may be different than the current live source.
	liveRune, err := gh.scopeResult.SourceTracker.LineAndColToRunePosition(lineNumber, colPosition, source, compilercommon.SourceMapCurrent)
	if err != nil {
		return []ContextOrAction{}, err
	}

	sourcePosition := source.PositionForRunePosition(liveRune, gh.scopeResult.SourceTracker)
	return gh.GetPositionalActions(sourcePosition)
}

// GetPositionalActions returns all asynchronous code actions for the given source position. Unlike GetContextActions, the
// returns items must *all* be actions, and should be displayed in a selector menu, rather than inline in the code.
func (gh Handle) GetPositionalActions(sourcePosition compilercommon.SourcePosition) ([]ContextOrAction, error) {
	updatedPosition, err := gh.scopeResult.SourceTracker.GetPositionOffset(sourcePosition, packageloader.CurrentFilePosition)
	if err != nil {
		return []ContextOrAction{}, err
	}

	// Create a set of async actions for every import statement in the file that comes from VCS.
	source := sourcePosition.Source()
	module, found := gh.scopeResult.Graph.SourceGraph().FindModuleBySource(source)
	if !found {
		return []ContextOrAction{}, fmt.Errorf("Invalid or non-SRG source file")
	}

	imports := module.GetImports()
	actions := make([]ContextOrAction, 0, len(imports))

	for _, srgImport := range imports {
		// Make sure the import has a valid source range.
		sourceRange, hasSourceRange := srgImport.SourceRange()
		if !hasSourceRange {
			continue
		}

		// Make sure the import's range contains the position.
		containsPosition, err := sourceRange.ContainsPosition(updatedPosition)
		if err != nil || !containsPosition {
			continue
		}

		// Parse the import's source and its source kind from the import string.
		importSource, kind, err := srgImport.ParsedSource()
		if err != nil {
			continue
		}

		// Make sure this import is pulling from VCS.
		if kind != parser.ParsedImportTypeVCS {
			continue
		}

		// Parse the VCS path to determine if this import is a HEAD reference, a branch/commit, or a tag.
		parsed, err := vcs.ParseVCSPath(importSource)
		if err != nil {
			log.Printf("Got error when parsing VCS path `%s`: %v", importSource, err)
			continue
		}

		// Inspect the VCS source to determine the commit information and tags available.
		inspectInfo, err := gh.getInspectInfo(importSource, source)
		if err != nil {
			log.Printf("Got error when inspecting VCS path `%s`: %v", importSource, err)
			continue
		}

		// Add an upgrade/update command (if applicable)
		updateTitle := "Update"
		updateVersion, err := compilerutil.SemanticUpdateVersion(parsed.Tag(), inspectInfo.Tags, compilerutil.UpdateVersionMinor)
		if err != nil {
			updateTitle = "Upgrade"
			updateVersion, err = compilerutil.SemanticUpdateVersion("0.0.0", inspectInfo.Tags, compilerutil.UpdateVersionMajor)
		}

		if err == nil && updateVersion != "" && updateVersion != "0.0.0" {
			actions = append(actions, ContextOrAction{
				Range:  sourceRange,
				Title:  fmt.Sprintf("%s import to %s", updateTitle, updateVersion),
				Action: FreezeImport,
				ActionParams: map[string]interface{}{
					"source":        string(source),
					"import-source": parsed.URL(),
					"tag":           updateVersion,
				},
			})
		}

		// Add a freeze/unfreeze command.
		if parsed.Tag() != "" || parsed.BranchOrCommit() != "" {
			actions = append(actions, ContextOrAction{
				Range:  sourceRange,
				Title:  "Unfreeze import to HEAD",
				Action: UnfreezeImport,
				ActionParams: map[string]interface{}{
					"source":        string(source),
					"import-source": parsed.URL(),
				},
			})
		} else if len(inspectInfo.CommitSHA) >= shortSHALength {
			actions = append(actions, ContextOrAction{
				Range:  sourceRange,
				Title:  "Freeze import at " + inspectInfo.CommitSHA[0:shortSHALength],
				Action: FreezeImport,
				ActionParams: map[string]interface{}{
					"source":        string(source),
					"import-source": parsed.URL(),
					"commit":        inspectInfo.CommitSHA,
				},
			})
		}
	}

	return actions, nil
}

// GetContextActions returns all context actions for the given source file.
func (gh Handle) GetContextActions(source compilercommon.InputSource) ([]CodeContextOrAction, error) {
	// Create a CodeContextOrAction for every import statement in the file that comes from VCS.
	module, found := gh.scopeResult.Graph.SourceGraph().FindModuleBySource(source)
	if !found {
		return []CodeContextOrAction{}, fmt.Errorf("Invalid or non-SRG source file")
	}

	imports := module.GetImports()
	cca := make([]CodeContextOrAction, 0, len(imports))

	for _, srgImport := range imports {
		// Make sure the import has a valid source range.
		sourceRange, hasSourceRange := srgImport.SourceRange()
		if !hasSourceRange {
			continue
		}

		// Parse the import's source and its source kind from the import string.
		importSource, kind, err := srgImport.ParsedSource()
		if err != nil {
			continue
		}

		// Make sure this import is pulling from VCS.
		if kind != parser.ParsedImportTypeVCS {
			continue
		}

		// Add Context to show the current commit SHA.
		cca = append(cca, CodeContextOrAction{
			Range: sourceRange,
			Resolve: func() (ContextOrAction, bool) {
				inspectInfo, err := gh.getInspectInfo(importSource, source)
				if err != nil {
					return ContextOrAction{}, false
				}

				if len(inspectInfo.CommitSHA) < shortSHALength {
					return ContextOrAction{}, false
				}

				return ContextOrAction{
					Range:        sourceRange,
					Title:        inspectInfo.CommitSHA[0:shortSHALength] + " (" + inspectInfo.Engine + ")",
					Action:       NoAction,
					ActionParams: map[string]interface{}{},
				}, true
			},
		})
	}

	return cca, nil
}

// getInspectInfo returns the VCS inspection information for the given import source, referenced under the given
// source file.
func (gh Handle) getInspectInfo(importSource string, source compilercommon.InputSource) (vcs.InspectInfo, error) {
	// First check the cache, indexed by the import source.
	cached, hasCached := gh.importInspectCache.Get(importSource)
	if hasCached {
		return cached.(vcs.InspectInfo), nil
	}

	// If not found, retrieve it via the VCS engine.
	dirPath := gh.groker.pathLoader.VCSPackageDirectory(gh.groker.entrypoint)
	inspectInfo, err, _ := vcs.PerformVCSCheckoutAndInspect(importSource, dirPath, vcs.VCSAlwaysUseCache)
	if err != nil {
		return vcs.InspectInfo{}, nil
	}

	gh.importInspectCache.Set(importSource, inspectInfo)
	return inspectInfo, nil
}
