package cmd

import (
	"github.com/Vilsol/go-mlog/cli"
	"github.com/Vilsol/go-mlog/runtime"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(executeCmd)
}

var executeCmd = &cobra.Command{
	Use:   "execute [flags] <program>",
	Short: "Execute MLOG",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO Output streaming
		objects, err := cli.ConstructObjectsFromConfig()

		if err != nil {
			return err
		}

		if err := runtime.ExecuteMLOGFile(args[0], objects); err != nil {
			return err
		}

		return nil
	},
}
