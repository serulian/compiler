// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package main is the serulian compiler.
package main

import (
	"fmt"
	"os"

	"github.com/serulian/compiler/builder"
	"github.com/serulian/compiler/developer"
	"github.com/serulian/compiler/formatter"
	"github.com/serulian/compiler/tester"
	"github.com/spf13/cobra"

	_ "github.com/serulian/compiler/tester/karma"
)

var (
	vcsDevelopmentDirectories []string
	debug                     bool
	addr                      string
)

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

			if !builder.BuildSource(args[0], debug, vcsDevelopmentDirectories...) {
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
			cmd.Usage()
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

			if !formatter.Format(args[0], debug) {
				os.Exit(-1)
			}
		},
	}

	var cmdFreeze = &cobra.Command{
		Use:   "freeze [source path] [vcs import]",
		Short: "Freezes the specified VCS imports in all Serulian files at the given path",
		Long:  `Modifies all imports of the given VCS libraries to refer to the SHA of the current HEAD commit`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected source path")
				os.Exit(-1)
			}

			if len(args) < 2 {
				fmt.Println("Expected one or more VCS imports")
				os.Exit(-1)
			}

			if !formatter.Freeze(args[0], args[1:len(args)], vcsDevelopmentDirectories, debug) {
				os.Exit(-1)
			}
		},
	}

	var cmdUnfreeze = &cobra.Command{
		Use:   "unfreeze [source path] [vcs import]",
		Short: "Unfreezes the specified VCS imports in all Serulian files at the given path",
		Long:  `Modifies all imports of the given VCS libraries to refer to HEAD`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected source path")
				os.Exit(-1)
			}

			if len(args) < 2 {
				fmt.Println("Expected one or more VCS imports")
				os.Exit(-1)
			}

			if !formatter.Unfreeze(args[0], args[1:len(args)], vcsDevelopmentDirectories, debug) {
				os.Exit(-1)
			}
		},
	}

	var cmdImports = &cobra.Command{
		Use:   "imports",
		Short: "Commands for modifying imports",
		Long:  "Commands for modifying imports",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
			os.Exit(1)
		},
	}

	// Register command-specific flags.
	cmdBuild.PersistentFlags().StringSliceVar(&vcsDevelopmentDirectories, "vcs-dev-dir", []string{},
		"If specified, VCS packages without specification will be first checked against this path")

	cmdDevelop.PersistentFlags().StringSliceVar(&vcsDevelopmentDirectories, "vcs-dev-dir", []string{},
		"If specified, VCS packages without specification will be first checked against this path")
	cmdDevelop.PersistentFlags().StringVar(&addr, "addr", ":8080", "The address at which the development code will be served")

	cmdImports.AddCommand(cmdFreeze)
	cmdImports.AddCommand(cmdUnfreeze)
	cmdImports.PersistentFlags().StringSliceVar(&vcsDevelopmentDirectories, "vcs-dev-dir", []string{},
		"If specified, VCS packages without specification will be first checked against this path")

	// Decorate the test commands.
	tester.DecorateRunners(cmdTest)

	// Register the root command.
	var rootCmd = &cobra.Command{
		Use:   "serulian",
		Short: "Serulian",
		Long:  "Serulian: A web and mobile development language",
	}

	rootCmd.AddCommand(cmdBuild)
	rootCmd.AddCommand(cmdDevelop)
	rootCmd.AddCommand(cmdFormat)
	rootCmd.AddCommand(cmdImports)
	rootCmd.AddCommand(cmdTest)
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "If set to true, Serulian will print debug logs")
	rootCmd.Execute()
}
