// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package main is the serulian compiler.
package main

import (
	"fmt"
	"os"

	"github.com/serulian/compiler/builder"
	"github.com/spf13/cobra"
)

var (
	vcsDevelopmentDirectories []string
	debug                     bool
)

func main() {
	var cmdBuild = &cobra.Command{
		Use:   "build [entrypoint source file]",
		Short: "Builds a Serulian project",
		Long:  `Builds a Serulian project, starting at the given root source file`,
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

	cmdBuild.PersistentFlags().BoolVar(&debug, "debug", false, "If set to true, Serulian will print debug logs during compilation")
	cmdBuild.PersistentFlags().StringSliceVar(&vcsDevelopmentDirectories, "vcs-dev-dir", []string{},
		"If specified, VCS packages without specification will be first checked against this path")

	var rootCmd = &cobra.Command{Use: "serulian"}
	rootCmd.AddCommand(cmdBuild)
	rootCmd.Execute()
}
