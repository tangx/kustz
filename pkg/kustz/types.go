package kustz

import (
	"os"

	"github.com/tangx/kustz/pkg/kubeutils"
)

const (
	APIVersion = "kustz/v1"
)

type Config struct {
	Metadata `json:",inline"`

	Name       string    `json:"name"`
	Namespace  string    `json:"namespace"`
	Service    Service   `json:"service"`
	Ingress    Ingress   `json:"ingress,omitempty"`
	ConfigMaps Generator `json:"configmaps,omitempty"`
	Secrets    Generator `json:"secrets,omitempty"`
	DNS        *DNS      `json:"dns,omitempty"`
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

	if kz.Metadata.APIVersion == "" {
		kz.Metadata.APIVersion = APIVersion
	}

	return kz
}

type Service struct {
	Name             string            `json:"name"`
	Image            string            `json:"image"`
	Replicas         int32             `json:"replicas"`
	Ports            []string          `json:"ports"`
	Envs             ServiceEnvs       `json:"envs,omitempty"`
	Resources        map[string]string `json:"resources,omitempty"`
	Probes           ContainerProbes   `json:"probes,omitempty"`
	ImagePullSecrets []string          `json:"imagePullSecrets,omitempty"`
	ImagePullPolicy  string            `json:"imagePullPolicy,omitempty"`
}

type ServiceEnvs struct {
	Pairs      map[string]string `json:"pairs,omitempty"`
	Files      []string          `json:"files,omitempty"`
	Secrets    []string          `json:"secrets,omitempty"`
	ConfigMaps []string          `json:"configmaps,omitempty"`
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

type ContainerProbes struct {
	Liveness  *ContainerProbe `json:"liveness,omitempty"`
	Readiness *ContainerProbe `json:"readiness,omitempty"`
	Startup   *ContainerProbe `json:"startup,omitempty"`
}

type ContainerProbe struct {
	ProbeHandler                  `json:",inline"`
	InitialDelaySeconds           int32  `json:"initialDelaySeconds,omitempty"`
	TimeoutSeconds                int32  `json:"timeoutSeconds,omitempty"`
	PeriodSeconds                 int32  `json:"periodSeconds,omitempty"`
	SuccessThreshold              int32  `json:"successThreshold,omitempty"`
	FailureThreshold              int32  `json:"failureThreshold,omitempty"`
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`
}

type ProbeHandler struct {
	Action  string            `json:"action,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}
