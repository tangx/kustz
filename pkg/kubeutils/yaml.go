package kubeutils

import (
	"os"

	"sigs.k8s.io/yaml"
)

func YAMLMarshal(v any) ([]byte, error) {
	b, err := yaml.Marshal(v)
	return b, err
}

func YAMLUnmarshal(data []byte, v interface{}) error {

	return yaml.Unmarshal(data, v)
}

func WriteYamlFile(name string, data any) error {
	b, err := YAMLMarshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(name, b, os.ModePerm)
}
