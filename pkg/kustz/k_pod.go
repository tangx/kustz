package kustz

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (kz *Config) KubePod() corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: kz.CommonLabels(),
		},
		Spec: corev1.PodSpec{
			Containers:       kz.KubeContainer(),
			ImagePullSecrets: toImagePullSecrets(kz.Service.ImagePullSecrets),
			DNSConfig:        toPodDNSConfig(kz.DNS),
			DNSPolicy:        toDNSPolicy(kz.DNS),
		},
	}
}

func toImagePullSecrets(secrets []string) []corev1.LocalObjectReference {
	if len(secrets) == 0 {
		return nil
	}

	objs := []corev1.LocalObjectReference{}
	for _, s := range secrets {
		objs = append(objs, corev1.LocalObjectReference{
			Name: s,
		})
	}

	return objs
}

func toPodDNSConfig(dns *DNS) *corev1.PodDNSConfig {
	if dns == nil {
		return nil
	}
	if dns.Config == nil {
		return nil
	}
	return &corev1.PodDNSConfig{
		Nameservers: dns.Config.Nameservers,
		Searches:    dns.Config.Searches,
		Options:     dns.Config.PodDNSConfigOptions(),
	}
}

func toDNSPolicy(dns *DNS) corev1.DNSPolicy {
	if dns == nil {
		// return v1.DNSNone
		return ""
	}

	return dns.DNSPolicy()
}
