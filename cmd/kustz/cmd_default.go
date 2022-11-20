package main

import (
	"github.com/spf13/cobra"
	"github.com/tangx/kustz/pkg/kustz"
)

var cmdDefault = &cobra.Command{
	Use:  "default",
	Long: "default",
	Run: func(cmd *cobra.Command, args []string) {
		kustz.Default()
	},
}

func init() {
	rootCmd.AddCommand(cmdDefault)
}
