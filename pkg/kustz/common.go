package kustz

func (kz *Config) CommonLabels() map[string]string {
	return map[string]string{
		"app": kz.Name,
	}
}
