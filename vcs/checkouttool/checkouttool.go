package main

import (
	"fmt"

	"github.com/Serulian/compiler/vcs"
	"github.com/spf13/cobra"
)

func main() {
	var path string
	var cmdRun = &cobra.Command{
		Use:   "checkouttool --path [path]",
		Short: "Small tool for testing VCS checkout",
		Long:  `A small tool for testing VCS checkout against a path`,
		Run: func(cmd *cobra.Command, args []string) {
			if path == "" {
				fmt.Println("Expected path")
				return
			}

			fmt.Printf("Performing checkout for path: %v\n", path)

			pkgCacheDirectory := ".pkg"
			packagePath, err, warning := vcs.PerformVCSCheckout(path, pkgCacheDirectory, vcs.VCSFollowNormalCacheRules)

			fmt.Printf("Path: %s\n", packagePath)
			fmt.Printf("Error: %v\n", err)
			fmt.Printf("Warning: %s\n", warning)
		},
	}

	cmdRun.Flags().StringVar(&path, "path", "", "The path to checkout")
	cmdRun.Execute()
}
