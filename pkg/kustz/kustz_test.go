package kustz

import (
	"fmt"
	"testing"

	"github.com/tangx/kustz/pkg/kubeutils"
)

func Test_Kustz(t *testing.T) {

	kz := NewKustzFromConfig("./kustz.yml")

	// dep := kz.KubeDeployment()
	// output(dep)

	svc := kz.KubeService()
	output(svc)
}

func output(v any) {
	b, err := kubeutils.YAMLMarshal(v)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", b)
}
