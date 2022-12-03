package kustz

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tangx/kustz/pkg/tokube"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
)

func (kz *Config) KubeContainer() []corev1.Container {
	if kz.Service.Name == "" {
		kz.Service.Name = kz.Name
	}

	probes := kz.Service.Probes

	c := corev1.Container{
		Name:            kz.Service.Name,
		Image:           kz.Service.Image,
		Env:             kz.kubeContainerEnv(),
		EnvFrom:         kz.kubeContainerEnvFrom(),
		Resources:       kz.kubeContainerResources(),
		LivenessProbe:   probes.kubeProbe(probes.Liveness),
		ReadinessProbe:  probes.kubeProbe(probes.Readiness),
		StartupProbe:    probes.kubeProbe(probes.Startup),
		ImagePullPolicy: tokube.ImagePullPolicy(kz.Service.ImagePullPolicy),
	}

	return []corev1.Container{c}
}

// kubeContainerEnvFrom 定义 configmap 或 secret 数据容器变量
// https://kubernetes.io/docs/concepts/configuration/secret/
func (kz *Config) kubeContainerEnvFrom() []corev1.EnvFromSource {

	sources := []corev1.EnvFromSource{}

	// value = config-name:true
	for _, value := range kz.Service.Envs.Secrets {
		sources = append(sources, tokube.ParseEnvFromSource(value, "secret"))
	}

	for _, value := range kz.Service.Envs.ConfigMaps {
		sources = append(sources, tokube.ParseEnvFromSource(value, "configmap"))
	}

	return sources
}

func (kz *Config) kubeContainerEnv() []corev1.EnvVar {
	pairs := make(map[string]string, 0)

	for _, file := range kz.Service.Envs.Files {
		b, err := os.ReadFile(file)
		if err != nil {
			logrus.Fatalf("read env file failed: %v", err)
		}
		err = yaml.Unmarshal(b, &pairs)
		if err != nil {
			logrus.Fatalf("unmarshal env file failed: %v", err)
		}
	}

	for k, v := range kz.Service.Envs.Pairs {
		pairs[k] = v
	}

	return tokube.ContainerEnv(pairs)
}

// kubeContainerResources 返回容器资源申请
// cpu and memory request: https://kubernetes.io/zh-cn/docs/concepts/configuration/manage-resources-containers/
// nvidia gpu request: https://help.aliyun.com/document_detail/94800.html
func (kz *Config) kubeContainerResources() corev1.ResourceRequirements {
	return tokube.ContainerResources(kz.Service.Resources)
}

func (cps ContainerProbes) kubeProbe(cp *ContainerProbe) *corev1.Probe {
	if cp == nil {
		return nil
	}
	return cp.kubeProbe()
}

// kubeProbe return Kube Probe without handler
func (cp *ContainerProbe) kubeProbe() *corev1.Probe {
	handler := tokube.ProbeHandler(cp.Action, cp.Headers)
	return &corev1.Probe{
		ProbeHandler:                  handler,
		InitialDelaySeconds:           cp.InitialDelaySeconds,
		TimeoutSeconds:                cp.TimeoutSeconds,
		PeriodSeconds:                 cp.PeriodSeconds,
		SuccessThreshold:              cp.SuccessThreshold,
		FailureThreshold:              cp.FailureThreshold,
		TerminationGracePeriodSeconds: cp.TerminationGracePeriodSeconds,
	}
}
