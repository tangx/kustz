package tokube

import (
	"net/url"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ProbeHandler action
//
// http(s)://:8080/healthy
// tcp://:8080
// cat /tmp/healthy
func ProbeHandler(action string, headers map[string]string) corev1.ProbeHandler {
	if strings.HasPrefix(action, "tcp://") {
		return toTCPProbeHandler(action)
	}

	if strings.HasPrefix(action, "http://") || strings.HasSuffix(action, "https://") {
		return toHTTPProbeHandler(action, headers)
	}

	return toExecProbeHandler(action)
}

func toHTTPProbeHandler(action string, headers map[string]string) corev1.ProbeHandler {

	ur, err := url.Parse(action)
	if err != nil {
		panic(err)
	}

	handler := corev1.ProbeHandler{
		HTTPGet: &corev1.HTTPGetAction{
			Scheme:      corev1.URIScheme(ur.Scheme),
			Host:        ur.Hostname(),
			Port:        intstr.Parse(ur.Port()),
			Path:        ur.Path,
			HTTPHeaders: toHTTPHeaders(headers),
		},
	}
	return handler
}

func toHTTPHeaders(headers map[string]string) []corev1.HTTPHeader {
	if len(headers) == 0 {
		return nil
	}

	hh := []corev1.HTTPHeader{}
	for k, v := range headers {
		hh = append(hh, corev1.HTTPHeader{
			Name:  k,
			Value: v,
		})
	}
	return hh
}

func toTCPProbeHandler(action string) corev1.ProbeHandler {
	ur, err := url.Parse(action)
	if err != nil {
		panic(err)
	}

	handler := corev1.ProbeHandler{
		TCPSocket: &corev1.TCPSocketAction{
			Host: ur.Hostname(),
			Port: intstr.Parse(ur.Port()),
		},
	}

	return handler
}

func toExecProbeHandler(action string) corev1.ProbeHandler {
	return corev1.ProbeHandler{
		Exec: &corev1.ExecAction{
			Command: []string{"sh", "-c", action},
		},
	}
}
