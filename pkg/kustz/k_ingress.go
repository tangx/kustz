package kustz

import (
	"net/url"
	"strings"

	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (kz *Config) KubeIngress() *netv1.Ingress {

	rules, tlss := ParseIngreseRulesFromStrings(kz.Ingress.Rules, kz.Name)
	ing := &netv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        kz.Name,
			Labels:      kz.CommonLabels(),
			Annotations: kz.Ingress.Annotations,
		},
		Spec: netv1.IngressSpec{
			Rules: rules,
			TLS:   tlss,
		},
	}

	return ing
}

func ParseIngreseRulesFromStrings(values []string, defaultService string) ([]netv1.IngressRule, []netv1.IngressTLS) {

	rules := []netv1.IngressRule{}
	tlss := []netv1.IngressTLS{}
	for _, value := range values {
		ing := NewIngressRuleFromString(value)
		if ing == nil {
			continue
		}

		if ing.Service == "" {
			ing.Service = defaultService
		}

		rules = append(rules, ing.KubeIngressRule())
		if tls := ing.KubeIngressTLS(); tls != nil {
			tlss = append(tlss, *tls)
		}
	}
	return rules, tlss
}

type IngressRuleString struct {
	Host      string
	Path      string
	PathType  netv1.PathType
	TLSSecret string
	Service   string
}

func NewIngressRuleFromString(value string) *IngressRuleString {

	ur, err := url.Parse(value)
	if err != nil {
		return nil
	}

	// ex: /api/*
	path := ur.Path
	typ := netv1.PathTypeExact
	if strings.HasSuffix(path, "*") {
		path = strings.TrimSuffix(path, "*")
		typ = netv1.PathTypePrefix
	}

	ing := &IngressRuleString{
		Host:      ur.Hostname(),
		Path:      path,
		PathType:  typ,
		TLSSecret: ur.Query().Get("tls"),
		Service:   ur.Query().Get("svc"),
	}

	return ing
}

func (ir *IngressRuleString) KubeIngressTLS() *netv1.IngressTLS {
	if ir.TLSSecret == "" {
		return nil
	}

	return &netv1.IngressTLS{
		Hosts: []string{
			ir.Host,
		},
		SecretName: ir.TLSSecret,
	}
}

func (ir *IngressRuleString) KubeIngressRule() netv1.IngressRule {

	ing := netv1.IngressRule{
		Host: ir.Host,
		IngressRuleValue: netv1.IngressRuleValue{
			HTTP: &netv1.HTTPIngressRuleValue{
				Paths: []netv1.HTTPIngressPath{
					{
						Path:     ir.Path,
						PathType: &ir.PathType,
						Backend:  ir.toKubeIngressBackend(),
					},
				},
			},
		},
	}

	return ing
}

func (ir *IngressRuleString) toKubeIngressBackend() netv1.IngressBackend {

	// srv-webapp-demo[:8080]
	svc := ir.Service
	port := int32(80)

	parts := strings.Split(svc, ":")
	if len(parts) == 2 {
		svc = parts[0]
		port = StringToInt32(parts[1])
	}

	return netv1.IngressBackend{
		Service: &netv1.IngressServiceBackend{
			Name: svc,
			Port: netv1.ServiceBackendPort{
				Number: int32(port),
			},
		},
	}
}
