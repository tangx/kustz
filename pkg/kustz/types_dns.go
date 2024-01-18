package kustz

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
)

type DNS struct {
	Config *DNSConfig `json:"config"`
	Policy string     `json:"policy"`
}

func (dns *DNS) DNSPolicy() corev1.DNSPolicy {
	if dns.Config != nil && dns.Policy == "" {
		return corev1.DNSNone
	}
	return corev1.DNSPolicy(dns.Policy)
}

type DNSConfig struct {
	Nameservers []string `json:"nameservers,omitempty"`
	Searches    []string `json:"searches,omitempty"`
	Options     []string `json:"options,omitempty"`
}

func (dc *DNSConfig) PodDNSConfigOptions() []corev1.PodDNSConfigOption {
	ret := make([]corev1.PodDNSConfigOption, 0)
	for i := range dc.Options {
		part := strings.Split(dc.Options[i], ":")

		switch len(part) {
		case 2:
			ret = append(ret, corev1.PodDNSConfigOption{
				Name:  part[0],
				Value: &part[1],
			})
		case 1:
			ret = append(ret, corev1.PodDNSConfigOption{
				Name: part[0],
			})
		default:
			continue
		}
	}
	return ret
}
