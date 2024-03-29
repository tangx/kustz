package kustz

import (
	"fmt"
	"testing"

	"github.com/tangx/kustz/pkg/kubeutils"
)

var (
	kz = NewKustzFromConfig("./kustz.yml")
)

func Test_YamlMarshal(t *testing.T) {
	b, err := kubeutils.YamlPkgMarshal(kz)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", b)
}

func Test_KustzDeployment(t *testing.T) {
	dep := kz.KubeDeployment()
	output(dep)
}

func Test_KustzService(t *testing.T) {
	svc := kz.KubeService()
	output(svc)
}

func Test_KustzIngress(t *testing.T) {
	ing := kz.KubeIngress()
	output(ing)
}

func Test_KustzKustomize(t *testing.T) {
	kust := kz.Kustomization()
	output(kust)
}

func output(v any) {
	b, err := kubeutils.YamlSigMarshal(v)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", b)
}
