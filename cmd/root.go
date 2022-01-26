package cmd

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	buildVersion = ""
	buildDate    = ""
	buildCommit  = ""
)

var rootCmd = &cobra.Command{
	Use:   "go-mlog",
	Short: "golang to mlog transpiler",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.SetEnvPrefix("gomlog")
		viper.AutomaticEnv()

		_ = viper.ReadInConfig()

		level, err := zerolog.ParseLevel(viper.GetString("log"))

		if err != nil {
			panic(err)
		}

		zerolog.SetGlobalLevel(level)

		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	},
}

func Execute(version string, date string, commit string) {
	buildVersion = version
	buildDate = date
	buildCommit = commit
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.PersistentFlags().String("log", "info", "The log level to output")
	rootCmd.PersistentFlags().Bool("colors", false, "Force log output with colors")

	rootCmd.PersistentFlags().Bool("numbers", false, "Output line numbers")
	rootCmd.PersistentFlags().Bool("comments", false, "Output comments")
	rootCmd.PersistentFlags().Int("comment-offset", 60, "Comment offset from line start")
	rootCmd.PersistentFlags().String("stacked", "", "Use a provided memory cell/bank as a stack")
	rootCmd.PersistentFlags().Bool("source", false, "Output source code after comment")

	rootCmd.PersistentFlags().String("output", "", "Output file. Outputs to stdout if unspecified")

	rootCmd.PersistentFlags().Bool("metrics", false, "Output source metrics after execution")

	_ = viper.BindPFlag("log", rootCmd.PersistentFlags().Lookup("log"))
	_ = viper.BindPFlag("colors", rootCmd.PersistentFlags().Lookup("colors"))

	_ = viper.BindPFlag("numbers", rootCmd.PersistentFlags().Lookup("numbers"))
	_ = viper.BindPFlag("comments", rootCmd.PersistentFlags().Lookup("comments"))
	_ = viper.BindPFlag("comment-offset", rootCmd.PersistentFlags().Lookup("comment-offset"))
	_ = viper.BindPFlag("stacked", rootCmd.PersistentFlags().Lookup("stacked"))
	_ = viper.BindPFlag("source", rootCmd.PersistentFlags().Lookup("source"))

	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))

	_ = viper.BindPFlag("metrics", rootCmd.PersistentFlags().Lookup("metrics"))
}
