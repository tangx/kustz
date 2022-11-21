package kustz

import (
	_ "embed"
	"fmt"

	"github.com/tangx/kustz/pkg/kubeutils"
)

//go:embed kustz.yml
var defaultConfig string

func DefaultConfig() {
	fmt.Println(defaultConfig)
}

func (kz *Config) RenderAll() error {
	dep := kz.KubeDeployment()
	ing := kz.KubeIngress()
	svc := kz.KubeService()
	kust := kz.Kustomization()

	datas := []struct {
		name string
		data any
	}{
		{name: FileDeployment, data: dep},
		{name: FileService, data: svc},
		{name: FileIngress, data: ing},
		{name: FileKustomization, data: kust},
	}

	for _, obj := range datas {
		err := kubeutils.WriteYamlFile(obj.name, obj.data)
		if err != nil {
			return err
		}
	}

	return nil
}

const (
	FileDeployment    = "deployment.yml"
	FileIngress       = "ingress.yml"
	FileService       = "service.yml"
	FileKustomization = "kustomization.yml"
)
