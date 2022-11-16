package kustz

import (
	corev1 "k8s.io/api/core/v1"
)

func (kz *Config) KubeContainer() []corev1.Container {
	if kz.Service.Name == "" {
		kz.Service.Name = kz.Name
	}

	c := corev1.Container{
		Name:  kz.Service.Name,
		Image: kz.Service.Image,
	}

	return []corev1.Container{c}
}
