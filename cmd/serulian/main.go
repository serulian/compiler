// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package main is the serulian compiler.
package main

import (
	"fmt"
	"os"
	runtime "runtime/debug"

	goprofile "github.com/pkg/profile"
	"github.com/serulian/compiler/builder"
	"github.com/serulian/compiler/developer"
	"github.com/serulian/compiler/formatter"
	"github.com/serulian/compiler/integration"
	"github.com/serulian/compiler/packagetools"
	"github.com/serulian/compiler/tester"
	"github.com/serulian/compiler/version"

	"github.com/spf13/cobra"

	_ "github.com/serulian/compiler/tester/karma"
)

var (
	vcsDevelopmentDirectories []string
	debug                     bool
	profile                   bool
	addr                      string
	verbose                   bool
	revisionNote              string
	upgrade                   bool
	yes                       bool
)

func disableGC() {
	// Disables GC.
	runtime.SetGCPercent(-1)
}

func main() {
	var cmdBuild = &cobra.Command{
		Use:   "build [entrypoint source file]",
		Short: "Builds a Serulian project",
		Long:  `Builds a Serulian project, starting at the given root source file.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Expected entrypoint source file")
				os.Exit(-1)
			}

			disableGC()

			if profile {
				defer goprofile.Start(goprofile.CPUProfile).Stop()
			}

			if !builder.BuildSource(args[0], debug, vcsDevelopmentDirectories...) && !profile {
				os.Exit(-1)
			}
		},
	}

	var cmdDevelop = &cobra.Command{
		Use:   "develop [entrypoint source file]",
		Short: "Starts development mode of a Serulian project",
		Long:  `Starts a webserver that automatically compiles on refresh, starting at the given root source file.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Expected entrypoint source file")
				os.Exit(-1)
			}

			if !developer.Run(addr, args[0], debug, vcsDevelopmentDirectories) {
				os.Exit(-1)
			}
		},
	}

	var cmdTest = &cobra.Command{
		Use:   "test",
		Short: "Runs the tests defined at the given source path",
		Long:  "Runs the tests found in any *_test.seru files at the given source path",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		},
	}

	var cmdFormat = &cobra.Command{
		Use:   "format [source path]",
		Short: "Formats all Serulian files at the given path",
		Long:  `Formats all Serulian files (.seru) found at the given path to the defined formatting.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Expected source path")
				os.Exit(-1)
			}

			if !formatter.Format(args[0], upgrade, debug) {
				os.Exit(-1)
			}
		},
	}

	var cmdFreeze = &cobra.Command{
		Use:   "freeze [source path] [vcs import]",
		Short: "Freezes imports",
		Long:  `Modifies all imports of the given VCS libraries to refer to the SHA of the current HEAD commit`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected source path")
				os.Exit(-1)
			}

			if len(args) < 2 {
				fmt.Println("Expected one or more VCS import patterns")
				os.Exit(-1)
			}

			if !formatter.Freeze(args[0], args[1:len(args)], vcsDevelopmentDirectories, debug) {
				os.Exit(-1)
			}
		},
	}

	var cmdUnfreeze = &cobra.Command{
		Use:   "unfreeze [source path] [vcs import]",
		Short: "Unfreezes imports",
		Long:  `Modifies all imports of the given VCS libraries to refer to HEAD`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected source path")
				os.Exit(-1)
			}

			if len(args) < 2 {
				fmt.Println("Expected one or more VCS import patterns")
				os.Exit(-1)
			}

			if !formatter.Unfreeze(args[0], args[1:len(args)], vcsDevelopmentDirectories, debug) {
				os.Exit(-1)
			}
		},
	}

	var cmdUpgrade = &cobra.Command{
		Use:   "upgrade [source path] [vcs import]",
		Short: "Upgrades imports",
		Long:  `Freezes the specified VCS imports in all Serulian files at the given path at the latest *stable* version`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected source path")
				os.Exit(-1)
			}

			if len(args) < 2 {
				fmt.Println("Expected one or more VCS import patterns")
				os.Exit(-1)
			}

			if !formatter.Upgrade(args[0], args[1:len(args)], vcsDevelopmentDirectories, debug) {
				os.Exit(-1)
			}
		},
	}

	var cmdUpdate = &cobra.Command{
		Use:   "update [source path] [vcs import]",
		Short: "Updates imports",
		Long:  `Changes the specified imports to the latest version *compatible* version, as per semvar`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected source path")
				os.Exit(-1)
			}

			if len(args) < 2 {
				fmt.Println("Expected one or more VCS import patterns")
				os.Exit(-1)
			}

			if !formatter.Update(args[0], args[1:len(args)], vcsDevelopmentDirectories, debug) {
				os.Exit(-1)
			}
		},
	}

	var cmdImports = &cobra.Command{
		Use:   "imports",
		Short: "Commands for modifying imports",
		Long:  "Commands for modifying imports",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		},
	}

	var cmdIntegrations = &cobra.Command{
		Use:   "integrations",
		Short: "Commands for working with integrations",
		Long:  "Commands for listing, installing, uninstalling and upgrading integrations",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		},
	}

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Displays the version of the toolkit",
		Long:  "Displays the version of the toolkit",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Serulian Toolkit\n\n")
			fmt.Printf("Toolkit Version: %s\n", version.DescriptiveVersion())
			fmt.Printf("Toolkit SHA: %s\n", version.GitSHA)
		},
	}

	var cmdPackage = &cobra.Command{
		Use:   "package",
		Short: "Serulian package management tools",
		Long:  `Tools and commands for managing updating and versioning of Serulian packages`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		},
	}

	var cmdDiff = &cobra.Command{
		Use:   "diff [package path] (compare-version)",
		Short: "Performs a diff of this package against latest tagged",
		Long:  `Performs a diff between the current version of this package and the latest tagged (or that specified)`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected package path")
				os.Exit(-1)
			}

			compareVersion := ""
			if len(args) > 1 {
				compareVersion = args[1]
			}

			packagetools.OutputDiff(args[0], compareVersion, verbose, debug, vcsDevelopmentDirectories...)
		},
	}

	var cmdRev = &cobra.Command{
		Use:   "rev [package path] (compare-version)",
		Short: "Tags this package with a semantic version computed against latest tagged",
		Long:  `Tags this package with a semantic version computed against latest tagged (or that specified)`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected package path")
				os.Exit(-1)
			}

			compareVersion := ""
			if len(args) > 1 {
				compareVersion = args[1]
			}

			packagetools.Revise(args[0], compareVersion, revisionNote, verbose, debug, vcsDevelopmentDirectories...)
		},
	}

	var cmdListIntegrations = &cobra.Command{
		Use:   "list",
		Short: "Lists all installed integrations",
		Long:  `Lists all Serulian integrations installed`,
		Run: func(cmd *cobra.Command, args []string) {
			if !integration.ListIntegrations() {
				os.Exit(1)
			}
		},
	}

	var cmdDescribeIntegration = &cobra.Command{
		Use:   "describe [id]",
		Short: "Describes an integration",
		Long:  `Describes an installed Serulian integration`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected ID")
				os.Exit(-1)
			}

			if !integration.DescribeIntegration(args[0]) {
				os.Exit(1)
			}
		},
	}

	var cmdInstallIntegration = &cobra.Command{
		Use:   "install [integration repository URL]",
		Short: "Installs an integration",
		Long:  `Installs a Serulian integration from its VCS repository URL`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected url")
				os.Exit(-1)
			}

			if !integration.InstallIntegration(args[0], debug) {
				os.Exit(1)
			}
		},
	}

	var cmdUninstallIntegration = &cobra.Command{
		Use:   "uninstall [id]",
		Short: "Uninstalls an integration",
		Long:  `Uninstalls an installed Serulian integration`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected ID")
				os.Exit(-1)
			}

			if !integration.UninstallIntegration(args[0], yes) {
				os.Exit(1)
			}
		},
	}

	var cmdUpgradeIntegration = &cobra.Command{
		Use:   "upgrade [id]",
		Short: "Upgrades an integration",
		Long:  `Upgrades an installed Serulian integration`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected ID")
				os.Exit(-1)
			}

			if !integration.UpgradeIntegration(args[0], yes, debug) {
				os.Exit(1)
			}
		},
	}

	cmdIntegrations.AddCommand(cmdListIntegrations)
	cmdIntegrations.AddCommand(cmdDescribeIntegration)
	cmdIntegrations.AddCommand(cmdInstallIntegration)
	cmdIntegrations.AddCommand(cmdUninstallIntegration)
	cmdIntegrations.AddCommand(cmdUpgradeIntegration)

	cmdPackage.AddCommand(cmdDiff)
	cmdPackage.AddCommand(cmdRev)

	// Register command-specific flags.
	cmdBuild.PersistentFlags().StringSliceVar(&vcsDevelopmentDirectories, "vcs-dev-dir", []string{},
		"If specified, VCS packages without specification will be first checked against this path")

	cmdDevelop.PersistentFlags().StringSliceVar(&vcsDevelopmentDirectories, "vcs-dev-dir", []string{},
		"If specified, VCS packages without specification will be first checked against this path")
	cmdDevelop.PersistentFlags().StringVar(&addr, "addr", ":8080", "The address at which the development code will be served")

	cmdTest.PersistentFlags().StringSliceVar(&vcsDevelopmentDirectories, "vcs-dev-dir", []string{},
		"If specified, VCS packages without specification will be first checked against this path")

	cmdFormat.PersistentFlags().BoolVarP(&upgrade, "upgrade", "u", false,
		"If true, older forms of source code syntax are supported for parsing and formatting")

	cmdImports.AddCommand(cmdFreeze)
	cmdImports.AddCommand(cmdUnfreeze)
	cmdImports.AddCommand(cmdUpgrade)
	cmdImports.AddCommand(cmdUpdate)
	cmdImports.PersistentFlags().StringSliceVar(&vcsDevelopmentDirectories, "vcs-dev-dir", []string{},
		"If specified, VCS packages without specification will be first checked against this path")

	cmdDiff.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"Whether to show full diff output")

	cmdRev.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"Whether to show full diff output")

	cmdRev.PersistentFlags().StringVarP(&revisionNote, "revision-note", "n", "",
		"If specified, the note to use when tagging the new version")

	cmdUpgradeIntegration.PersistentFlags().BoolVarP(&yes, "yes", "y", false,
		"If true, the prompt will be skipped")

	cmdUninstallIntegration.PersistentFlags().BoolVarP(&yes, "yes", "y", false,
		"If true, the prompt will be skipped")

	// Decorate the test commands.
	tester.DecorateRunners(cmdTest, &vcsDevelopmentDirectories)

	// Register the root command.
	var rootCmd = &cobra.Command{
		Use:   "serulian",
		Short: fmt.Sprintf("Serulian %s", version.DescriptiveVersion()),
		Long:  fmt.Sprintf("Serulian %s: A web and mobile development language", version.DescriptiveVersion()),
	}

	rootCmd.AddCommand(cmdBuild)
	rootCmd.AddCommand(cmdDevelop)
	rootCmd.AddCommand(cmdTest)
	rootCmd.AddCommand(cmdFormat)
	rootCmd.AddCommand(cmdImports)
	rootCmd.AddCommand(cmdPackage)
	rootCmd.AddCommand(cmdIntegrations)
	rootCmd.AddCommand(cmdVersion)
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "If set to true, Serulian will print debug logs")
	rootCmd.PersistentFlags().BoolVar(&profile, "profile", false, "If set to true, Serulian will be profiled")
	rootCmd.Execute()
}
