package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dlx",
	Short: "dlx is a downloader",
}

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dlx-0.0.1")
		},
	}
	downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "download command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("download stuff")
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd, downloadCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
