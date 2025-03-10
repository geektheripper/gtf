package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/geektheripper/go-gutils/git/git_utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logger = log.New(os.Stderr)

var rootCmd = &cobra.Command{
	Use:   "gtf",
	Short: "GeekTR's command line terraform utilities",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func init() {
	viper.SetEnvPrefix("GTF")

	defaultRepo, _ := git_utils.FindGitRoot(".")
	rootCmd.PersistentFlags().StringP("repo", "r", defaultRepo, "the remote repository to manage")
	viper.BindEnv("repo")

	viper.BindPFlags(rootCmd.PersistentFlags())
}
