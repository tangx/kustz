package kustz

import (
	"fmt"
	"testing"

	"github.com/tangx/kustz/pkg/kubeutils"
)

var (
	kz = NewKustzFromConfig("./kustz.yml")
)

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

func output(v any) {
	b, err := kubeutils.YAMLMarshal(v)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", b)
}
