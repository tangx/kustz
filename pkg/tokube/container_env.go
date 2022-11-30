package tokube

import (
	"strconv"
	"strings"

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

func ParseEnvFromSource(value string, kind string) corev1.EnvFromSource {
	opt := false
	var err error
	parts := strings.Split(value, ":")
	if len(parts) == 2 {
		value = parts[0]
		opt, err = strconv.ParseBool(parts[1])
		if err != nil {
			opt = false
		}
	}

	switch kind {
	case "secret":
		return corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: value,
				},
				Optional: &opt,
			},
		}
	case "configmap":
		fallthrough
	default:
		return corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: value,
				},
				Optional: &opt,
			},
		}
	}
}
