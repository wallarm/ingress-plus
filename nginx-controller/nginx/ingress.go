package nginx

import (
	api_v1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
)

// IngressEx holds an Ingress along with Secrets and Endpoints of the services
// that are referenced in this Ingress
type IngressEx struct {
	Ingress    *extensions.Ingress
	TLSSecrets map[string]*api_v1.Secret
	JWTKey     *api_v1.Secret
	Endpoints  map[string][]string
}

type MergeableIngresses struct {
	Master  *IngressEx
	Minions []*IngressEx
}

var masterBlacklist = []string{
	"nginx.org/rewrites",
	"nginx.org/ssl-services",
	"nginx.org/websocket-services",
	"nginx.com/sticky-cookie-services",
}

var minionBlacklist = []string{
	"nginx.org/proxy-hide-headers",
	"nginx.org/proxy-pass-headers",
	"nginx.org/redirect-to-https",
	"ingress.kubernetes.io/ssl-redirect",
	"nginx.org/hsts",
	"nginx.org/hsts-max-age",
	"nginx.org/hsts-include-subdomains",
	"nginx.org/server-tokens",
	"nginx.org/listen-ports",
	"nginx.org/listen-ports-ssl",
	"nginx.com/jwt-key",
	"nginx.com/jwt-realm",
	"nginx.com/jwt-token",
	"nginx.com/jwt-login-url",
}
