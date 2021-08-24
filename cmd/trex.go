package cmd

import (
	"github.com/Vilsol/go-mlog/cli"
	"github.com/Vilsol/go-mlog/runtime"
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(trexCmd)
}

var trexCmd = &cobra.Command{
	Use:   "trex [flags] <program>",
	Short: "Transpile Go to MLOG and execute it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := transpiler.GolangToMLOGFile(args[0], transpiler.Options{
			Numbers:       viper.GetBool("numbers"),
			Comments:      viper.GetBool("comments"),
			CommentOffset: viper.GetInt("comment-offset"),
			Stacked:       viper.GetString("stacked"),
			Source:        viper.GetBool("source"),
		})

		if err != nil {
			return err
		}

		objects, err := cli.ConstructObjectsFromConfig()

		if err != nil {
			return err
		}

		// TODO Output streaming
		if _, err := runtime.ExecuteMLOG(result, objects); err != nil {
			return err
		}

		return nil
	},
}
