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
	"github.com/spf13/cobra"
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

	cmdFormat.PersistentFlags().BoolVar(&debug, "debug", false, "If set to true, Serulian will print debug logs during formatting")

	cmdBuild.PersistentFlags().BoolVar(&debug, "debug", false, "If set to true, Serulian will print debug logs during compilation")
	cmdBuild.PersistentFlags().StringSliceVar(&vcsDevelopmentDirectories, "vcs-dev-dir", []string{},
		"If specified, VCS packages without specification will be first checked against this path")

	cmdDevelop.PersistentFlags().BoolVar(&debug, "debug", false, "If set to true, Serulian will print debug logs during compilation")
	cmdDevelop.PersistentFlags().StringSliceVar(&vcsDevelopmentDirectories, "vcs-dev-dir", []string{},
		"If specified, VCS packages without specification will be first checked against this path")
	cmdDevelop.PersistentFlags().StringVar(&addr, "addr", ":8080", "The address at which the development code will be served")

	var rootCmd = &cobra.Command{Use: "serulian"}
	rootCmd.AddCommand(cmdBuild)
	rootCmd.AddCommand(cmdDevelop)
	rootCmd.AddCommand(cmdFormat)
	rootCmd.Execute()
}
