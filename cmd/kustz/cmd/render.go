package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tangx/kustz/pkg/kustz"
)

var cmdRender = &cobra.Command{
	Use:   "render",
	Short: "读取 kustz 配置， 生成 kustomize 所需文件",
	Run: func(cmd *cobra.Command, args []string) {

		kz := kustz.NewKustzFromConfig(config)
		kz.RenderAll()
	},
}

func init() {
	rootCmd.AddCommand(cmdRender)
	cmdRender.Flags().StringVarP(&config, "config", "c", "kustz.yml", "kustz config")
}

var config string
