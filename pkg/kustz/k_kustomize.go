package kustz

import (
	"github.com/tangx/kustz/pkg/tokust"
	"sigs.k8s.io/kustomize/v3/pkg/types"
)

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
		ConfigMapGenerator: kz.KustConfigMapGenerator(),
		SecretGenerator:    kz.KustSecretGenerator(),
	}

	return k
}

func (kz *Config) KustConfigMapGenerator() []types.ConfigMapArgs {

	args := []types.ConfigMapArgs{}

	for _, arg := range kz.ConfigMaps.Literals {
		arg := tokust.ConfigMapArgs_Literals(arg.Name, arg.Files)
		args = append(args, arg)
	}

	for _, arg := range kz.ConfigMaps.Files {
		arg := tokust.ConfigMapArgs_Files(arg.Name, arg.Files)
		args = append(args, arg)
	}
	for _, arg := range kz.ConfigMaps.Envs {
		arg := tokust.ConfigMapArgs_Env(arg.Name, arg.Files)
		args = append(args, arg)
	}

	return args
}

func (kz *Config) KustSecretGenerator() []types.SecretArgs {

	args := []types.SecretArgs{}

	for _, arg := range kz.Secrets.Literals {
		arg := tokust.SecretArgs_Liternals(arg.Name, arg.Files, arg.Type)
		args = append(args, arg)
	}

	for _, arg := range kz.Secrets.Files {
		arg := tokust.SecretArgs_Files(arg.Name, arg.Files, arg.Type)
		args = append(args, arg)
	}
	for _, arg := range kz.Secrets.Envs {
		arg := tokust.SecretArgs_Env(arg.Name, arg.Files, arg.Type)
		args = append(args, arg)
	}

	return args
}
