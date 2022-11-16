package kustz

import (
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (kz *Config) KubeDeployment() *appv1.Deployment {
	return &appv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kz.Name,
			Namespace: kz.Namespace,
			Labels:    CommonLabels(*kz),
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &kz.Service.Replicas,
			Template: kz.KubePod(),
			Selector: &metav1.LabelSelector{
				MatchLabels: CommonLabels(*kz),
			},
		},
	}
}
