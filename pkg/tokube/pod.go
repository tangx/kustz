package tokube

import corev1 "k8s.io/api/core/v1"

func ImagePullSecrets(secrets []string) []corev1.LocalObjectReference {
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
