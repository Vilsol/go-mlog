package cmd

import (
	"fmt"
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

func init() {
	rootCmd.AddCommand(transpileCmd)
}

var transpileCmd = &cobra.Command{
	Use:   "transpile [flags] <program>",
	Short: "Transpile Go to MLOG",
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

		if output := viper.GetString("output"); output != "" {
			if err := ioutil.WriteFile(output, []byte(result), 0644); err != nil {
				return err
			}
		} else {
			fmt.Println(result)
		}

		return nil
	},
}
