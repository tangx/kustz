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
		Name:  kz.Service.Name,
		Image: kz.Service.Image,
		Env:   kz.kubeContainerEnv(),
	}

	return []corev1.Container{c}
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

// func (kz *Config) kubeContainerEnv_Error() []corev1.EnvVar {
// 	envs := []corev1.EnvVar{}

// 	envs = append(envs, tokube.ContainerEnv(kz.Service.Envs.Pairs)...)

// 	for _, file := range kz.Service.Envs.Files {
// 		b, err := os.ReadFile(file)
// 		if err != nil {
// 			logrus.Fatalf("read env file failed: %v", err)
// 		}
// 		mm := make(map[string]string, 0)
// 		err = yaml.Unmarshal(b, &mm)
// 		if err != nil {
// 			logrus.Fatalf("unmarshal env file failed: %v", err)
// 		}

// 		envs = append(envs, tokube.ContainerEnv(mm)...)
// 	}
// 	return envs
// }
