package kustz

import (
	"os"

	"github.com/tangx/kustz/pkg/kubeutils"
)

type Config struct {
	Name       string    `json:"name"`
	Namespace  string    `json:"namespace"`
	Service    Service   `json:"service"`
	Ingress    Ingress   `json:"ingress"`
	ConfigMaps Generator `json:"configmaps"`
	Secrets    Generator `json:"secrets"`
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
	Name     string      `json:"name"`
	Image    string      `json:"image"`
	Replicas int32       `json:"replicas"`
	Ports    []string    `json:"ports"`
	Envs     ServiceEnvs `json:"envs,omitempty"`
}

type ServiceEnvs struct {
	Pairs map[string]string `json:"pairs,omitempty"`
	Files []string          `json:"files,omitempty"`
}

type Ingress struct {
	Rules       []string          `json:"rules"`
	Annotations map[string]string `json:"annotations"`
}

// Generator 定义数据源种类
type Generator struct {
	Literals []GeneratorArgs `json:"literals,omitempty"`
	Envs     []GeneratorArgs `json:"envs,omitempty"`
	Files    []GeneratorArgs `json:"files,omitempty"`
}

// GeneratorArgs 定义数据源类型参数
type GeneratorArgs struct {
	Name  string   `json:"name,omitempty"`
	Files []string `json:"files,omitempty"`
	Type  string   `json:"type,omitempty"`
}
