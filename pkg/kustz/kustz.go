package kustz

type Kustz struct {
	Namespace string  `yaml:"namespace"`
	Service   Service `yaml:"service"`
}

type Service struct {
	Name     string `yaml:"name"`
	Image    string `yaml:"image"`
	Replicas int    `yaml:"replicas"`
}
