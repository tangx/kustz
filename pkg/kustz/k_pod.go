package kustz

import (
	"github.com/tangx/kustz/pkg/tokube"
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
			ImagePullSecrets: tokube.ImagePullSecrets(kz.Service.ImagePullSecrets),
		},
	}
}
