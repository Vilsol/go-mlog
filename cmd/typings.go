package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Vilsol/go-mlog/checker"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(typingsCmd)
}

var typingsCmd = &cobra.Command{
	Use:   "typings",
	Short: "Output typings as JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		result := checker.GetSerializablePackages()
		marshal, _ := json.Marshal(result)
		fmt.Println(string(marshal))
		return nil
	},
}
