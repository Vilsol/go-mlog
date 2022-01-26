package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current go-mlog version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%s-%s-%s", buildVersion, buildCommit, buildDate)
		return nil
	},
}
