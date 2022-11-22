package tokube

import (
	corev1 "k8s.io/api/core/v1"
)

func ContainerEnv(pairs map[string]string) []corev1.EnvVar {
	envs := []corev1.EnvVar{}
	for k, v := range pairs {
		envs = append(envs, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	return envs
}
