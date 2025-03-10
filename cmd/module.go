package cmd

import (
	"github.com/spf13/cobra"
)

var moduleCmd = &cobra.Command{
	Use:     "module",
	Aliases: []string{"m"},
	Short:   "terraform module management",
}

func init() {
	rootCmd.AddCommand(moduleCmd)
}
