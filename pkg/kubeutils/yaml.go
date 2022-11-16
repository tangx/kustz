package kubeutils

import "sigs.k8s.io/yaml"

func YAMLMarshal(v any) ([]byte, error) {
	b, err := yaml.Marshal(v)
	return b, err
}

func YAMLUnmarshal(data []byte, v interface{}) error {

	return yaml.Unmarshal(data, v)
}
