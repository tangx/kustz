package tokube

import (
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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

// ContainerResources return kube ResourceRequirements
// resources:
//
//	cpu: 10m/20m
//	memory: 0Mi/20Mi
func ContainerResources(res map[string]string) corev1.ResourceRequirements {
	// 如果资源为空， 直接返回
	if len(res) == 0 {
		return corev1.ResourceRequirements{}
	}

	limits := corev1.ResourceList{}
	requests := corev1.ResourceList{}

	for k, v := range res {
		name := corev1.ResourceName(k)
		req, limit := toResourceQuantity(v)
		limits[name] = limit
		requests[name] = req
	}
	return corev1.ResourceRequirements{
		Limits:   limits,
		Requests: requests,
	}
}

// toResourceQuantity return both request and limit resource.Quantity
// value should obey: request[/limit]
// below is legal:
//
//	for cpu: 10m/20m
//	for mem: 0Mi/20Mi
func toResourceQuantity(value string) (request resource.Quantity, limit resource.Quantity) {

	re, li := "", ""
	parts := strings.Split(value, "/")
	if len(parts) == 1 {
		re = value
		li = value
	}
	if len(parts) == 2 {
		re = parts[0]
		li = parts[1]
	}

	request = resource.MustParse(re)
	limit = resource.MustParse(li)

	return
}
