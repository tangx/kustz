package tokust

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
	"sigs.k8s.io/kustomize/v3/pkg/types"
)

func marshalYaml(file string, out any) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, out)
	return err
}

func mapToSlice(in map[string]string) []string {
	out := []string{}
	for k, v := range in {
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}

	return out
}

func generatorArgs_literals(name string, files []string) types.GeneratorArgs {
	data := make(map[string]string, 0)
	for _, file := range files {
		err := marshalYaml(file, data)
		if err != nil {
			panic(err)
		}
	}

	sources := mapToSlice(data)
	g := types.GeneratorArgs{
		Name: name,
		DataSources: types.DataSources{
			LiteralSources: sources,
		},
	}

	return g
}

func generatorArgs_files(name string, files []string) types.GeneratorArgs {

	g := types.GeneratorArgs{
		Name: name,
		DataSources: types.DataSources{
			FileSources: files,
		},
	}

	return g
}

func generatorArgs_envs(name string, files []string) types.GeneratorArgs {

	g := types.GeneratorArgs{
		Name: name,
		DataSources: types.DataSources{
			EnvSources: files,
		},
	}

	return g
}

func ConfigMapArgs_Env(name string, files []string) types.ConfigMapArgs {
	return types.ConfigMapArgs{
		GeneratorArgs: generatorArgs_envs(name, files),
	}
}

func ConfigMapArgs_Files(name string, files []string) types.ConfigMapArgs {
	return types.ConfigMapArgs{
		GeneratorArgs: generatorArgs_files(name, files),
	}
}
func ConfigMapArgs_Literals(name string, files []string) types.ConfigMapArgs {
	return types.ConfigMapArgs{
		GeneratorArgs: generatorArgs_literals(name, files),
	}
}

func SecretArgs_Env(name string, files []string, typ string) types.SecretArgs {
	if typ == "" {
		typ = "Opaque"
	}
	return types.SecretArgs{
		GeneratorArgs: generatorArgs_envs(name, files),
		Type:          typ,
	}
}

func SecretArgs_Files(name string, files []string, typ string) types.SecretArgs {
	if typ == "" {
		typ = "Opaque"
	}
	return types.SecretArgs{
		GeneratorArgs: generatorArgs_files(name, files),
		Type:          typ,
	}
}

func SecretArgs_Liternals(name string, files []string, typ string) types.SecretArgs {
	if typ == "" {
		typ = "Opaque"
	}
	return types.SecretArgs{
		GeneratorArgs: generatorArgs_literals(name, files),
		Type:          typ,
	}
}
