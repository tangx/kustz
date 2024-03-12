package kubeutils

import (
	"os"

	pkgyaml "gopkg.in/yaml.v3"
	sigyaml "sigs.k8s.io/yaml"
)

func YamlSigMarshal(v any) ([]byte, error) {
	b, err := sigyaml.Marshal(v)
	return b, err
}

func YamlSigUnmarshal(data []byte, v interface{}) error {

	return sigyaml.Unmarshal(data, v)
}

func WriteYamlFile(name string, data any) error {
	b, err := YamlSigMarshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(name, b, os.ModePerm)
}

func YamlPkgMarshal(v any) ([]byte, error) {
	b, err := pkgyaml.Marshal(v)
	return b, err
}
func YamlPkgUnmarshal(data []byte, v interface{}) error {
	return pkgyaml.Unmarshal(data, v)
}
