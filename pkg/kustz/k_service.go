package kustz

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (kz *Config) KubeService() *corev1.Service {

	ports, typ := ParsePortStrings(kz.Service.Ports)

	svc := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   kz.Name,
			Labels: kz.CommonLabels(),
			// Namespace: kz.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: kz.CommonLabels(),
			Type:     typ,
			Ports:    ports,
		},
	}

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

type PortString struct {
	Port       int32
	TargetPort int32
	NodePort   int32
	Protocol   corev1.Protocol
	Type       corev1.ServiceType
}

// NewPortFromString parse port from string PortString
// port like
//
//	tcp://!10080:80:8080 => tcp/udp/sctp
//	8080:80
//	!18080:80:8080
func NewPortFromString(value string) PortString {
	port := &PortString{
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
func (p *PortString) KubeServicePort() corev1.ServicePort {

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

// toServiceClusterIP parse value from for ClusterIP
func (p *PortString) toServiceClusterIP(value string) {

	parts := strings.Split(value, ":")
	switch len(parts) {
	case 1:
		n := StringToInt32(parts[0])
		p.Port = n
		p.TargetPort = n
	case 2:
		p.Port = StringToInt32(parts[0])
		p.TargetPort = StringToInt32(parts[1])
	}

	p.Type = corev1.ServiceTypeClusterIP
}

// toServiceNodePort parse value from for NodePort
func (p *PortString) toServiceNodePort(value string) {

	value = strings.TrimPrefix(value, "!")
	parts := strings.Split(value, ":")
	switch len(parts) {
	case 1, 2:
		p.toServiceClusterIP(value)
	case 3:
		p.NodePort = StringToInt32(parts[0])
		p.Port = StringToInt32(parts[1])
		p.TargetPort = StringToInt32(parts[2])
	}

	p.Type = corev1.ServiceTypeNodePort
}
