package tokube

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func ImagePullPolicy(s string) corev1.PullPolicy {

	switch strings.ToLower(s) {
	case "always":
		return corev1.PullAlways
	case "never":
		return corev1.PullNever
	case "ifnotpresent":
		return corev1.PullIfNotPresent
	}

	return ""
}
