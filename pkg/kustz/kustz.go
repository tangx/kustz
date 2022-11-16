package kustz

import (
	"os"

	"github.com/tangx/kustz/pkg/kubeutils"
)

type Config struct {
	Name      string  `json:"name"`
	Namespace string  `json:"namespace"`
	Service   Service `json:"service"`
}

func NewKustzFromConfig(cfg string) *Config {
	b, err := os.ReadFile(cfg)
	if err != nil {
		panic(err)
	}

	kz := &Config{}
	err = kubeutils.YAMLUnmarshal(b, kz)
	if err != nil {
		panic(err)
	}

	return kz
}

type Service struct {
	Name     string   `json:"name"`
	Image    string   `json:"image"`
	Replicas int32    `json:"replicas"`
	Ports    []string `json:"ports"`
}
