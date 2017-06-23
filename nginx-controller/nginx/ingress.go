package nginx

import (
	api_v1 "k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

// IngressEx holds an Ingress along with Secrets and Endpoints of the services
// that are referenced in this Ingress
type IngressEx struct {
	Ingress   *extensions.Ingress
	Secrets   map[string]*api_v1.Secret
	Endpoints map[string][]string
}
