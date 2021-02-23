package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
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

		level, err := log.ParseLevel(viper.GetString("log"))

		if err != nil {
			panic(err)
		}

		log.SetFormatter(&log.TextFormatter{
			ForceColors: viper.GetBool("colors"),
		})
		log.SetOutput(os.Stdout)
		log.SetLevel(level)
	},
}

func Execute() {
	// Allow running from explorer
	cobra.MousetrapHelpText = ""

	// Execute transpile command as default
	cmd, _, err := rootCmd.Find(os.Args[1:])
	if (len(os.Args) <= 1 || os.Args[1] != "help") && (err != nil || cmd == rootCmd) {
		args := append([]string{"transpile"}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

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

	_ = viper.BindPFlag("log", rootCmd.PersistentFlags().Lookup("log"))
	_ = viper.BindPFlag("colors", rootCmd.PersistentFlags().Lookup("colors"))

	_ = viper.BindPFlag("numbers", rootCmd.PersistentFlags().Lookup("numbers"))
	_ = viper.BindPFlag("comments", rootCmd.PersistentFlags().Lookup("comments"))
	_ = viper.BindPFlag("comment-offset", rootCmd.PersistentFlags().Lookup("comment-offset"))
	_ = viper.BindPFlag("stacked", rootCmd.PersistentFlags().Lookup("stacked"))
	_ = viper.BindPFlag("source", rootCmd.PersistentFlags().Lookup("source"))

	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
}
