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

	c := corev1.Container{
		Name:    kz.Service.Name,
		Image:   kz.Service.Image,
		Env:     kz.kubeContainerEnv(),
		EnvFrom: kz.kubeContainerEnvFrom(),
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
