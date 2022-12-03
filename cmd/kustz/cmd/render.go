package cmd

import (
	"github.com/go-jarvis/cobrautils"
	"github.com/spf13/cobra"
	"github.com/tangx/kustz/pkg/kustz"
)

var cmdRender = &cobra.Command{
	Use:   "render",
	Short: "读取 kustz 配置， 生成 kustomize 所需文件",
	Run: func(cmd *cobra.Command, args []string) {
		render()
	},
}

func init() {
	rootCmd.AddCommand(cmdRender)
	cobrautils.BindFlags(cmdRender, flags)
}

type KustzFlag struct {
	Config   string `flag:"config" usage:"kustz config" shorthand:"c"`
	Image    string `flag:"image" usage:"image name"`
	Replicas *int   `flag:"replicas" usage:"pod replicas number"`
}

var flags = &KustzFlag{
	Config: "kustz.yml",
}

func render() {
	kz := kustz.NewKustzFromConfig(flags.Config)

	if flags.Image != "" {
		kz.Service.Image = flags.Image
	}
	if flags.Replicas != nil {
		kz.Service.Replicas = int32(*flags.Replicas)
	}

	kz.RenderAll()
}
