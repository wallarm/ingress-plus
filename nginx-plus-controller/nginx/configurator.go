package nginx

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

const emptyHost = ""

// Configurator transforms an Ingress resource into NGINX Configuration
type Configurator struct {
	nginx    *NginxController
	nginxAPI *NginxAPIController
	config   *Config
	lock     sync.Mutex
}

// NewConfigurator creates a new Configurator
func NewConfigurator(nginx *NginxController, config *Config, nginxAPI *NginxAPIController) *Configurator {
	cnf := Configurator{
		nginx:    nginx,
		config:   config,
		nginxAPI: nginxAPI,
	}

	return &cnf
}

// AddOrUpdateIngress adds or updates NGINX configuration for an Ingress resource
func (cnf *Configurator) AddOrUpdateIngress(name string, ingEx *IngressEx) {
	cnf.lock.Lock()
	defer cnf.lock.Unlock()

	pems := cnf.updateCertificates(ingEx)
	nginxCfg := cnf.generateNginxCfg(ingEx, pems)
	cnf.nginx.AddOrUpdateIngress(name, nginxCfg)
	cnf.nginx.Reload()
	time.Sleep(500 * time.Millisecond)
	cnf.updateEndpoints(name, ingEx)
}

func (cnf *Configurator) updateCertificates(ingEx *IngressEx) map[string]string {
	pems := make(map[string]string)

	for _, tls := range ingEx.Ingress.Spec.TLS {
		secretName := tls.SecretName
		secret, exist := ingEx.Secrets[secretName]
		if !exist {
			continue
		}
		cert, ok := secret.Data[api.TLSCertKey]
		if !ok {
			glog.Warningf("Secret %v has no private key", secretName)
			continue
		}
		key, ok := secret.Data[api.TLSPrivateKeyKey]
		if !ok {
			glog.Warningf("Secret %v has no cert", secretName)
			continue
		}

		pemFileName := cnf.nginx.AddOrUpdateCertAndKey(secretName, string(cert), string(key))

		for _, host := range tls.Hosts {
			pems[host] = pemFileName
		}
		if len(tls.Hosts) == 0 {
			pems[emptyHost] = pemFileName
		}
	}

	return pems
}
func (cnf *Configurator) generateNginxCfg(ingEx *IngressEx, pems map[string]string) IngressNginxConfig {
	ingCfg := cnf.createConfig(ingEx)

	upstreams := make(map[string]Upstream)

	if ingEx.Ingress.Spec.Backend != nil {
		name := getNameForUpstream(ingEx.Ingress, emptyHost, ingEx.Ingress.Spec.Backend.ServiceName)
		upstream := cnf.createUpstream(name)
		upstreams[name] = upstream
	}

	var servers []Server

	for _, rule := range ingEx.Ingress.Spec.Rules {
		if rule.IngressRuleValue.HTTP == nil {
			continue
		}

		serverName := rule.Host

		statuzZone := rule.Host
		if rule.Host == emptyHost {
			statuzZone = ingEx.Ingress.Namespace + "-" + ingEx.Ingress.Name
			glog.Warningf("Host field of ingress rule in %v/%v is empty", ingEx.Ingress.Namespace, ingEx.Ingress.Name)
		}

		server := Server{Name: serverName, StatusZone: statuzZone}

		if pemFile, ok := pems[serverName]; ok {
			server.SSL = true
			server.SSLCertificate = pemFile
			server.SSLCertificateKey = pemFile
		}

		var locations []Location
		rootLocation := false

		for _, path := range rule.HTTP.Paths {
			upsName := getNameForUpstream(ingEx.Ingress, rule.Host, path.Backend.ServiceName)

			if _, exists := upstreams[upsName]; !exists {
				upstream := cnf.createUpstream(upsName)
				upstreams[upsName] = upstream
			}

			loc := createLocation(pathOrDefault(path.Path), upstreams[upsName], &ingCfg)
			locations = append(locations, loc)

			if loc.Path == "/" {
				rootLocation = true
			}
		}

		if rootLocation == false && ingEx.Ingress.Spec.Backend != nil {
			upsName := getNameForUpstream(ingEx.Ingress, emptyHost, ingEx.Ingress.Spec.Backend.ServiceName)
			loc := createLocation(pathOrDefault("/"), upstreams[upsName], &ingCfg)
			locations = append(locations, loc)
		}

		server.Locations = locations
		servers = append(servers, server)
	}

	if len(ingEx.Ingress.Spec.Rules) == 0 && ingEx.Ingress.Spec.Backend != nil {
		serverName := emptyHost
		statuzZone := ingEx.Ingress.Namespace + "-" + ingEx.Ingress.Name
		glog.Warningf("Host field of ingress rule in %v/%v is empty", ingEx.Ingress.Namespace, ingEx.Ingress.Name)

		server := Server{Name: serverName, StatusZone: statuzZone}

		if pemFile, ok := pems[emptyHost]; ok {
			server.SSL = true
			server.SSLCertificate = pemFile
			server.SSLCertificateKey = pemFile
		}

		var locations []Location

		upsName := getNameForUpstream(ingEx.Ingress, emptyHost, ingEx.Ingress.Spec.Backend.ServiceName)

		loc := createLocation(pathOrDefault("/"), upstreams[upsName], &ingCfg)
		locations = append(locations, loc)

		server.Locations = locations
		servers = append(servers, server)
	}

	return IngressNginxConfig{Upstreams: upstreamMapToSlice(upstreams), Servers: servers}
}

func (cnf *Configurator) createConfig(ingEx *IngressEx) Config {
	ingCfg := *cnf.config
	if proxyConnectTimeout, exists := ingEx.Ingress.Annotations["nginx.org/proxy-connect-timeout"]; exists {
		ingCfg.ProxyConnectTimeout = proxyConnectTimeout
	}
	if proxyReadTimeout, exists := ingEx.Ingress.Annotations["nginx.org/proxy-read-timeout"]; exists {
		ingCfg.ProxyReadTimeout = proxyReadTimeout
	}
	if clientMaxBodySize, exists := ingEx.Ingress.Annotations["nginx.org/client-max-body-size"]; exists {
		ingCfg.ClientMaxBodySize = clientMaxBodySize
	}

	return ingCfg
}

func createLocation(path string, upstream Upstream, cfg *Config) Location {
	loc := Location{
		Path:                path,
		Upstream:            upstream,
		ProxyConnectTimeout: cfg.ProxyConnectTimeout,
		ProxyReadTimeout:    cfg.ProxyReadTimeout,
		ClientMaxBodySize:   cfg.ClientMaxBodySize,
	}

	return loc
}

func (cnf *Configurator) createUpstream(name string) Upstream {
	return Upstream{Name: name}
}

func pathOrDefault(path string) string {
	if path == "" {
		return "/"
	}
	return path
}

func endpointsToUpstreamServers(endps api.Endpoints, servicePort int) []UpstreamServer {
	var upsServers []UpstreamServer
	for _, subset := range endps.Subsets {
		for _, port := range subset.Ports {
			if port.Port == servicePort {
				for _, address := range subset.Addresses {
					ups := UpstreamServer{Address: address.IP, Port: fmt.Sprintf("%v", servicePort)}
					upsServers = append(upsServers, ups)
				}
				break
			}
		}
	}

	return upsServers
}

func getNameForUpstream(ing *extensions.Ingress, host string, service string) string {
	return fmt.Sprintf("%v-%v-%v-%v", ing.Namespace, ing.Name, host, service)
}

func upstreamMapToSlice(upstreams map[string]Upstream) []Upstream {
	result := make([]Upstream, 0, len(upstreams))

	for _, ups := range upstreams {
		result = append(result, ups)
	}

	return result
}

// DeleteIngress deletes NGINX configuration for an Ingress resource
func (cnf *Configurator) DeleteIngress(name string) {
	cnf.lock.Lock()
	defer cnf.lock.Unlock()

	cnf.nginx.DeleteIngress(name)
	cnf.nginx.Reload()
}

// UpdateEndpoints updates endpoints in NGINX configuration for an Ingress resource
func (cnf *Configurator) UpdateEndpoints(name string, ingEx *IngressEx) {
	cnf.lock.Lock()
	defer cnf.lock.Unlock()

	cnf.updateEndpoints(name, ingEx)
}

func (cnf *Configurator) updateEndpoints(name string, ingEx *IngressEx) {
	if ingEx.Ingress.Spec.Backend != nil {
		name := getNameForUpstream(ingEx.Ingress, emptyHost, ingEx.Ingress.Spec.Backend.ServiceName)
		endps, exists := ingEx.Endpoints[ingEx.Ingress.Spec.Backend.ServiceName]
		if exists {
			endpoints := getEndpointsList(endps, ingEx.Ingress.Spec.Backend.ServicePort.IntValue())
			err := cnf.nginxAPI.UpdateServers(name, endpoints)
			if err != nil {
				glog.Warningf("Couldn't update the endponts for %v: %v", name, err)
			}
		}
	}
	for _, rule := range ingEx.Ingress.Spec.Rules {
		if rule.IngressRuleValue.HTTP == nil {
			continue
		}
		for _, path := range rule.HTTP.Paths {
			name := getNameForUpstream(ingEx.Ingress, rule.Host, path.Backend.ServiceName)
			endps, exists := ingEx.Endpoints[path.Backend.ServiceName]
			if exists {
				endpoints := getEndpointsList(endps, path.Backend.ServicePort.IntValue())
				err := cnf.nginxAPI.UpdateServers(name, endpoints)
				if err != nil {
					glog.Warningf("Couldn't update the endponts for %v: %v", name, err)
				}
			}
		}
	}
}

func getEndpointsList(endps *api.Endpoints, servicePort int) []string {
	var result []string

	for _, subset := range endps.Subsets {
		for _, port := range subset.Ports {
			if port.Port == servicePort {
				for _, address := range subset.Addresses {
					result = append(result, fmt.Sprintf("%v:%v", address.IP, servicePort))
				}
				break
			}
		}
	}

	return result
}

// UpdateConfig updates NGINX Configuration parameters
func (cnf *Configurator) UpdateConfig(config *Config) {
	cnf.lock.Lock()
	defer cnf.lock.Unlock()

	cnf.config = config
	mainCfg := &NginxMainConfig{
		ServerNamesHashBucketSize: config.MainServerNamesHashBucketSize,
		ServerNamesHashMaxSize:    config.MainServerNamesHashMaxSize,
	}

	cnf.nginx.UpdateMainConfigFile(mainCfg)
}
