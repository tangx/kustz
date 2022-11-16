package kustz

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (kz *Config) KubeService() *corev1.Service {

	svc := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kz.Name,
			Namespace: kz.Namespace,
			Labels:    kz.CommonLabels(),
		},
		Spec: corev1.ServiceSpec{
			Selector: kz.CommonLabels(),
		},
	}

	ports, typ := ParsePortStrings(kz.Service.Ports)
	svc.Spec.Type = typ
	svc.Spec.Ports = ports

	return svc
}

func ParsePortStrings(values []string) ([]corev1.ServicePort, corev1.ServiceType) {

	sps := []corev1.ServicePort{}
	typ := corev1.ServiceTypeClusterIP

	for _, value := range values {
		port := NewPortFromString(value)
		if port.Type != corev1.ServiceTypeClusterIP {
			typ = port.Type
		}
		sps = append(sps, port.KubeServicePort())

	}

	return sps, typ
}

type Port struct {
	Port       int32
	TargetPort int32
	NodePort   int32
	Protocol   corev1.Protocol
	Type       corev1.ServiceType
}

// NewPortFromString parse port from string to v1.ServicePort
// port like
//
//	tcp://!10080:80:8080 => tcp/udp/sctp
//	8080:80
//	!18080:80:8080
func NewPortFromString(value string) Port {
	port := &Port{
		Protocol: corev1.ProtocolTCP,
		Type:     corev1.ServiceTypeClusterIP,
	}
	parts := strings.Split(value, "://")
	if len(parts) == 2 {
		value = parts[1]

		proto := parts[0]
		switch strings.ToLower(proto) {
		case "udp":
			port.Protocol = corev1.ProtocolUDP
		case "sctp":
			port.Protocol = corev1.ProtocolSCTP
		default:
			port.Protocol = corev1.ProtocolTCP
		}
	}

	sign := value[0]
	switch sign {
	case '!':
		port.toServiceNodePort(value)
	default:
		port.toServiceClusterIP(value)
	}

	return *port
}

// KubeServicePort return a corev1.ServicePort
func (p *Port) KubeServicePort() corev1.ServicePort {

	sp := &corev1.ServicePort{
		Name:       fmt.Sprintf("%d-%d", p.Port, p.TargetPort),
		Port:       p.Port,
		TargetPort: intstr.FromInt(int(p.TargetPort)),
		Protocol:   p.Protocol,
	}

	if p.TargetPort != 0 {
		sp.NodePort = p.NodePort
	}
	return *sp
}

// toServiceClusterIP parse value to
func (p *Port) toServiceClusterIP(value string) {

	parts := strings.Split(value, ":")
	switch len(parts) {
	case 1:
		n := p.StringToInt32(parts[0])
		p.Port = n
		p.TargetPort = n
	case 2:
		p.Port = p.StringToInt32(parts[0])
		p.TargetPort = p.StringToInt32(parts[1])
	}

	p.Type = corev1.ServiceTypeClusterIP
}

func (p *Port) toServiceNodePort(value string) {

	value = strings.TrimPrefix(value, "!")
	parts := strings.Split(value, ":")
	switch len(parts) {
	case 1, 2:
		p.toServiceClusterIP(value)
	case 3:
		p.NodePort = p.StringToInt32(parts[0])
		p.Port = p.StringToInt32(parts[1])
		p.TargetPort = p.StringToInt32(parts[2])
	}

	p.Type = corev1.ServiceTypeNodePort
}

func (p *Port) StringToInt32(val string) int32 {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	if i > math.MaxInt32 || i < math.MinInt32 {
		return 0
	}
	return int32(i)
}
