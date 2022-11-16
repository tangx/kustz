package kustz

func CommonLabels(kz Config) map[string]string {
	return map[string]string{
		"app": kz.Name,
	}
}
