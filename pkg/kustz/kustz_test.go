package kustz

import (
	"fmt"
	"testing"

	"github.com/tangx/kustz/pkg/kubeutils"
)

func Test_Kustz(t *testing.T) {

	kz := NewKustzFromConfig("./kustz.yml")

	dep := kz.KubeDeployment()
	b, err := kubeutils.YAMLMarshal(dep)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", b)
}
