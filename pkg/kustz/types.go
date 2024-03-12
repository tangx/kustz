package kustz

import (
	"os"

	"github.com/tangx/kustz/pkg/kubeutils"
)

const (
	APIVersion = "kustz/v1"
)

type Config struct {
	Metadata `json:",inline" yaml:",inline"`

	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`

	Service    Service   `json:"service" yaml:"service"`
	ConfigMaps Generator `json:"configmaps,omitempty" yaml:"configmaps,omitempty"`
	Secrets    Generator `json:"secrets,omitempty" yaml:"secrets,omitempty"`

	Ingress Ingress `json:"ingress,omitempty" yaml:"ingress,omitempty"`
	DNS     *DNS    `json:"dns,omitempty" yaml:"dns,omitempty"`
}

func NewKustzFromConfig(cfg string) *Config {
	b, err := os.ReadFile(cfg)
	if err != nil {
		panic(err)
	}

	kz := &Config{}
	err = kubeutils.YamlSigUnmarshal(b, kz)
	if err != nil {
		panic(err)
	}

	if kz.Metadata.APIVersion == "" {
		kz.Metadata.APIVersion = APIVersion
	}

	return kz
}

type Service struct {
	Name             string            `json:"name" yaml:"name"`
	Image            string            `json:"image" yaml:"image"`
	Replicas         int32             `json:"replicas" yaml:"replicas"`
	Ports            []string          `json:"ports" yaml:"ports"`
	Envs             ServiceEnvs       `json:"envs,omitempty" yaml:"envs,omitempty"`
	Resources        map[string]string `json:"resources,omitempty" yaml:"resources,omitempty"`
	Probes           ContainerProbes   `json:"probes,omitempty" yaml:"probes,omitempty"`
	ImagePullSecrets []string          `json:"imagePullSecrets,omitempty" yaml:"imagePullSecrets,omitempty"`
	ImagePullPolicy  string            `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
}

type ServiceEnvs struct {
	Pairs      map[string]string `json:"pairs,omitempty" yaml:"pairs,omitempty"`
	Files      []string          `json:"files,omitempty" yaml:"files,omitempty"`
	Secrets    []string          `json:"secrets,omitempty" yaml:"secrets,omitempty"`
	ConfigMaps []string          `json:"configmaps,omitempty" yaml:"configmaps,omitempty"`
}

type Ingress struct {
	Rules       []string          `json:"rules" yaml:"rules"`
	Annotations map[string]string `json:"annotations" yaml:"annotations"`
}

// Generator 定义数据源种类
type Generator struct {
	Literals []GeneratorArgs `json:"literals,omitempty" yaml:"literals,omitempty"`
	Envs     []GeneratorArgs `json:"envs,omitempty" yaml:"envs,omitempty"`
	Files    []GeneratorArgs `json:"files,omitempty" yaml:"files,omitempty"`
}

// GeneratorArgs 定义数据源类型参数
type GeneratorArgs struct {
	Name  string   `json:"name,omitempty" yaml:"name,omitempty"`
	Files []string `json:"files,omitempty" yaml:"files,omitempty"`
	Type  string   `json:"type,omitempty" yaml:"type,omitempty"`
}

type ContainerProbes struct {
	Liveness  *ContainerProbe `json:"liveness,omitempty" yaml:"liveness,omitempty"`
	Readiness *ContainerProbe `json:"readiness,omitempty" yaml:"readiness,omitempty"`
	Startup   *ContainerProbe `json:"startup,omitempty" yaml:"startup,omitempty"`
}

type ContainerProbe struct {
	ProbeHandler                  `json:",inline" yaml:",inline"`
	InitialDelaySeconds           int32  `json:"initialDelaySeconds,omitempty" yaml:"initialDelaySeconds,omitempty"`
	TimeoutSeconds                int32  `json:"timeoutSeconds,omitempty" yaml:"timeoutSeconds,omitempty"`
	PeriodSeconds                 int32  `json:"periodSeconds,omitempty" yaml:"periodSeconds,omitempty"`
	SuccessThreshold              int32  `json:"successThreshold,omitempty" yaml:"successThreshold,omitempty"`
	FailureThreshold              int32  `json:"failureThreshold,omitempty" yaml:"failureThreshold,omitempty"`
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty" yaml:"terminationGracePeriodSeconds,omitempty"`
}

type ProbeHandler struct {
	Action  string            `json:"action,omitempty" yaml:"action,omitempty"`
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
}
