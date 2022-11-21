package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tangx/kustz/pkg/kustz"
)

var cmdDefault = &cobra.Command{
	Use:   "default",
	Short: "在屏幕上打印 kustz 默认配置",
	Run: func(cmd *cobra.Command, args []string) {
		kustz.DefaultConfig()
	},
}

func init() {
	rootCmd.AddCommand(cmdDefault)
}
