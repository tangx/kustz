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
		ConfigMapGenerator: kz.ConfigMaps.toConfigMapArgs(),
		SecretGenerator:    kz.Secrets.toSecretArgs(),
	}

	return k
}

// toConfigMapArgs 返回 ConfigMap 参数
func (genor *Generator) toConfigMapArgs() []types.ConfigMapArgs {

	args := []types.ConfigMapArgs{}

	for _, data := range genor.datas() {
		for _, garg := range data.gargs {
			arg := tokust.ConfigMapArgs(garg.Name, garg.Files, data.mode)
			args = append(args, arg)
		}
	}

	return args
}

// toSecretArgs 返回 Secret 参数
func (genor *Generator) toSecretArgs() []types.SecretArgs {

	args := []types.SecretArgs{}

	for _, data := range genor.datas() {
		for _, garg := range data.gargs {
			arg := tokust.SecretArgs(garg.Name, garg.Files, garg.Type, data.mode)
			args = append(args, arg)
		}
	}

	return args
}

type GeneratorArgsData struct {
	mode  tokust.GeneratorMode
	gargs []GeneratorArgs
}

// datas 整合生成器数据
func (genor *Generator) datas() []GeneratorArgsData {
	return []GeneratorArgsData{
		{mode: tokust.GeneratorMode_Envs, gargs: genor.Envs},
		{mode: tokust.GeneratorMode_Files, gargs: genor.Files},
		{mode: tokust.GeneratorMode_Literals, gargs: genor.Literals},
	}
}
