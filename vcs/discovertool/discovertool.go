package main

import (
	"fmt"

	"github.com/serulian/compiler/vcs"
	"github.com/spf13/cobra"
)

func main() {
	var url string
	var cmdRun = &cobra.Command{
		Use:   "discovertool --url [url]",
		Short: "Small tool for testing VCS discovery",
		Long:  `A small tool for testing VCS discovery against a URL`,
		Run: func(cmd *cobra.Command, args []string) {
			if url == "" {
				fmt.Println("Expected URL")
				return
			}

			fmt.Printf("Performing discovery for URL: %v\n", url)

			found, err := vcs.DiscoverVCSInformation(url)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			fmt.Printf("Discovery information: %v\n", found)
		},
	}

	cmdRun.Flags().StringVar(&url, "url", "", "The URL to discover")
	cmdRun.Execute()
}
