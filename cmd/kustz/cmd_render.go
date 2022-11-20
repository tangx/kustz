package main

import (
	"github.com/spf13/cobra"
	"github.com/tangx/kustz/pkg/kustz"
)

var cmdRender = &cobra.Command{
	Use:  "render",
	Long: "render",
	Run: func(cmd *cobra.Command, args []string) {

		kz := kustz.NewKustzFromConfig(config)
		kz.DumpAll()
	},
}

func init() {
	rootCmd.AddCommand(cmdRender)
	cmdRender.Flags().StringVarP(&config, "config", "c", "kustz.yml", "kustz config")
}

var config string
