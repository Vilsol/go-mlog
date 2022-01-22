package cmd

import (
	"fmt"
	"github.com/Vilsol/go-mlog/decompiler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

func init() {
	rootCmd.AddCommand(decompileCmd)
}

var decompileCmd = &cobra.Command{
	Use:   "decompile [flags] <program>",
	Short: "Decompile MLOG to Go",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := decompiler.MLOGToGolangFile(args[0])
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
