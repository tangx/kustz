package kustz

import "sigs.k8s.io/kustomize/v3/pkg/types"

func (kz *Config) Kustomization() types.Kustomization {
	k := types.Kustomization{
		TypeMeta: types.TypeMeta{
			Kind:       types.KustomizationKind,
			APIVersion: types.KustomizationVersion,
		},
		Namespace: kz.Namespace,
		Resources: []string{
			"deployment.yml",
			"ingress.yml",
			"service.yml",
		},
	}

	return k
}
