/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/golang/glog"

	"github.com/nginxinc/kubernetes-ingress/nginx-plus-controller/nginx"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/cache"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/watch"
)

// LoadBalancerController watches Kubernetes API and
// reconfigures NGINX via NGINXController when needed
type LoadBalancerController struct {
	client               *client.Client
	ingController        *framework.Controller
	svcController        *framework.Controller
	endpController       *framework.Controller
	cfgmController       *framework.Controller
	ingLister            StoreToIngressLister
	svcLister            cache.StoreToServiceLister
	endpLister           cache.StoreToEndpointsLister
	cfgmLister           StoreToConfigMapLister
	ingQueue             *taskQueue
	endpQueue            *taskQueue
	cfgmQueue            *taskQueue
	stopCh               chan struct{}
	cnf                  *nginx.Configurator
	watchNGINXConfigMaps bool
}

var keyFunc = framework.DeletionHandlingMetaNamespaceKeyFunc

// NewLoadBalancerController creates a controller
func NewLoadBalancerController(kubeClient *client.Client, resyncPeriod time.Duration, namespace string, cnf *nginx.Configurator, nginxConfigMaps string) (*LoadBalancerController, error) {
	lbc := LoadBalancerController{
		client: kubeClient,
		stopCh: make(chan struct{}),
		cnf:    cnf,
	}

	lbc.ingQueue = NewTaskQueue(lbc.syncIng)
	lbc.endpQueue = NewTaskQueue(lbc.syncEndp)

	ingHandlers := framework.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			addIng := obj.(*extensions.Ingress)
			glog.V(3).Infof("Adding Ingress: %v", addIng.Name)
			lbc.ingQueue.enqueue(obj)
		},
		DeleteFunc: func(obj interface{}) {
			remIng := obj.(*extensions.Ingress)
			glog.V(3).Infof("Removing Ingress: %v", remIng.Name)
			lbc.ingQueue.enqueue(obj)
		},
		UpdateFunc: func(old, cur interface{}) {
			if !reflect.DeepEqual(old, cur) {
				glog.V(3).Infof("Ingress %v changed, syncing",
					cur.(*extensions.Ingress).Name)
				lbc.ingQueue.enqueue(cur)
			}
		},
	}
	lbc.ingLister.Store, lbc.ingController = framework.NewInformer(
		&cache.ListWatch{
			ListFunc:  ingressListFunc(kubeClient, namespace),
			WatchFunc: ingressWatchFunc(kubeClient, namespace),
		},
		&extensions.Ingress{}, resyncPeriod, ingHandlers)

	svcHandlers := framework.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			addSvc := obj.(*api.Service)
			glog.V(3).Infof("Adding service: %v", addSvc.Name)
			lbc.enqueueIngressForService(obj)
		},
		DeleteFunc: func(obj interface{}) {
			remSvc := obj.(*api.Service)
			glog.V(3).Infof("Removing service: %v", remSvc.Name)
			lbc.enqueueIngressForService(obj)
		},
		UpdateFunc: func(old, cur interface{}) {
			if !reflect.DeepEqual(old, cur) {
				glog.V(3).Infof("Service %v changed, syncing",
					cur.(*api.Service).Name)
				lbc.enqueueIngressForService(cur)
			}
		},
	}
	lbc.svcLister.Store, lbc.svcController = framework.NewInformer(
		&cache.ListWatch{
			ListFunc:  serviceListFunc(kubeClient, namespace),
			WatchFunc: serviceWatchFunc(kubeClient, namespace),
		},
		&api.Service{}, resyncPeriod, svcHandlers)

	endpHandlers := framework.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			addEndp := obj.(*api.Endpoints)
			glog.V(3).Infof("Adding endpoints: %v", addEndp.Name)
			lbc.endpQueue.enqueue(obj)
		},
		DeleteFunc: func(obj interface{}) {
			remEndp := obj.(*api.Endpoints)
			glog.V(3).Infof("Removing endpoints: %v", remEndp.Name)
			lbc.endpQueue.enqueue(obj)
		},
		UpdateFunc: func(old, cur interface{}) {
			if !reflect.DeepEqual(old, cur) {
				glog.V(3).Infof("Endpoints %v changed, syncing",
					cur.(*api.Endpoints).Name)
				lbc.endpQueue.enqueue(cur)
			}
		},
	}
	lbc.endpLister.Store, lbc.endpController = framework.NewInformer(
		&cache.ListWatch{
			ListFunc:  endpointsListFunc(kubeClient, namespace),
			WatchFunc: endpointsWatchFunc(kubeClient, namespace),
		},
		&api.Endpoints{}, resyncPeriod, endpHandlers)

	if nginxConfigMaps != "" {
		nginxConfigMapsNS, nginxConfigMapsName, err := parseNGINXConfigMaps(nginxConfigMaps)
		if err != nil {
			glog.Warning(err)
		} else {
			lbc.watchNGINXConfigMaps = true
			lbc.cfgmQueue = NewTaskQueue(lbc.syncCfgm)

			cfgmHandlers := framework.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					cfgm := obj.(*api.ConfigMap)
					if cfgm.Name == nginxConfigMapsName {
						glog.V(3).Infof("Adding ConfigMap: %v", cfgm.Name)
						lbc.cfgmQueue.enqueue(obj)
					}
				},
				DeleteFunc: func(obj interface{}) {
					cfgm := obj.(*api.ConfigMap)
					if cfgm.Name == nginxConfigMapsName {
						glog.V(3).Infof("Removing ConfigMap: %v", cfgm.Name)
						lbc.cfgmQueue.enqueue(obj)
					}
				},
				UpdateFunc: func(old, cur interface{}) {
					if !reflect.DeepEqual(old, cur) {
						cfgm := cur.(*api.ConfigMap)
						if cfgm.Name == nginxConfigMapsName {
							glog.V(3).Infof("ConfigMap %v changed, syncing",
								cur.(*api.ConfigMap).Name)
							lbc.cfgmQueue.enqueue(cur)
						}
					}
				},
			}
			lbc.cfgmLister.Store, lbc.cfgmController = framework.NewInformer(
				&cache.ListWatch{
					ListFunc:  configMapsListFunc(kubeClient, nginxConfigMapsNS),
					WatchFunc: configMapsWatchFunc(kubeClient, nginxConfigMapsNS),
				},
				&api.ConfigMap{}, resyncPeriod, cfgmHandlers)
		}
	}

	return &lbc, nil
}

// Run starts the loadbalancer controller
func (lbc *LoadBalancerController) Run() {
	go lbc.ingController.Run(lbc.stopCh)
	go lbc.svcController.Run(lbc.stopCh)
	go lbc.endpController.Run(lbc.stopCh)
	go lbc.ingQueue.run(time.Second, lbc.stopCh)
	go lbc.endpQueue.run(time.Second, lbc.stopCh)
	if lbc.watchNGINXConfigMaps {
		go lbc.cfgmController.Run(lbc.stopCh)
		go lbc.cfgmQueue.run(time.Second, lbc.stopCh)
	}
	<-lbc.stopCh
}

func ingressListFunc(c *client.Client, ns string) func(api.ListOptions) (runtime.Object, error) {
	return func(opts api.ListOptions) (runtime.Object, error) {
		return c.Extensions().Ingress(ns).List(opts)
	}
}

func ingressWatchFunc(c *client.Client, ns string) func(options api.ListOptions) (watch.Interface, error) {
	return func(options api.ListOptions) (watch.Interface, error) {
		return c.Extensions().Ingress(ns).Watch(options)
	}
}

func serviceListFunc(c *client.Client, ns string) func(api.ListOptions) (runtime.Object, error) {
	return func(opts api.ListOptions) (runtime.Object, error) {
		return c.Services(ns).List(opts)
	}
}

func serviceWatchFunc(c *client.Client, ns string) func(options api.ListOptions) (watch.Interface, error) {
	return func(options api.ListOptions) (watch.Interface, error) {
		return c.Services(ns).Watch(options)
	}
}

func endpointsListFunc(c *client.Client, ns string) func(api.ListOptions) (runtime.Object, error) {
	return func(opts api.ListOptions) (runtime.Object, error) {
		return c.Endpoints(ns).List(opts)
	}
}

func endpointsWatchFunc(c *client.Client, ns string) func(options api.ListOptions) (watch.Interface, error) {
	return func(options api.ListOptions) (watch.Interface, error) {
		return c.Endpoints(ns).Watch(options)
	}
}

func configMapsListFunc(c *client.Client, ns string) func(api.ListOptions) (runtime.Object, error) {
	return func(opts api.ListOptions) (runtime.Object, error) {
		return c.ConfigMaps(ns).List(opts)
	}
}

func configMapsWatchFunc(c *client.Client, ns string) func(options api.ListOptions) (watch.Interface, error) {
	return func(options api.ListOptions) (watch.Interface, error) {
		return c.ConfigMaps(ns).Watch(options)
	}
}

func (lbc *LoadBalancerController) syncEndp(key string) {
	glog.V(3).Infof("Syncing endpoints %v", key)

	obj, endpExists, err := lbc.endpLister.Store.GetByKey(key)
	if err != nil {
		lbc.endpQueue.requeue(key, err)
		return
	}

	if endpExists {
		ings := lbc.getIngressForEndpoints(obj)

		for _, ing := range ings {
			ingEx := lbc.createIngress(&ing)
			glog.V(3).Infof("Updating Endponits for %v/%v", ing.Name, ing.Namespace)
			name := ing.Namespace + "-" + ing.Name
			lbc.cnf.UpdateEndpoints(name, &ingEx)
		}
	}

}

func (lbc *LoadBalancerController) syncCfgm(key string) {
	glog.V(3).Infof("Syncing configmap %v", key)

	obj, cfgmExists, err := lbc.cfgmLister.Store.GetByKey(key)
	if err != nil {
		lbc.cfgmQueue.requeue(key, err)
		return
	}
	cfg := nginx.NewDefaultConfig()

	if cfgmExists {
		cfgm := obj.(*api.ConfigMap)

		if proxyConnectTimeout, exists := cfgm.Data["proxy-connect-timeout"]; exists {
			cfg.ProxyConnectTimeout = proxyConnectTimeout
		}
		if proxyReadTimeout, exists := cfgm.Data["proxy-read-timeout"]; exists {
			cfg.ProxyReadTimeout = proxyReadTimeout
		}
	}
	lbc.cnf.UpdateConfig(cfg)

	ings, _ := lbc.ingLister.List()
	for _, ing := range ings.Items {
		lbc.ingQueue.enqueue(&ing)
	}
}

func (lbc *LoadBalancerController) syncIng(key string) {
	glog.V(3).Infof("Syncing %v", key)

	obj, ingExists, err := lbc.ingLister.Store.GetByKey(key)
	if err != nil {
		lbc.ingQueue.requeue(key, err)
		return
	}

	// defaut/some-ingress -> default-some-ingress
	name := strings.Replace(key, "/", "-", -1)

	if !ingExists {
		glog.V(2).Infof("Deleting Ingress: %v\n", key)
		lbc.cnf.DeleteIngress(name)
	} else {
		glog.V(2).Infof("Adding or Updating Ingress: %v\n", key)

		ing := obj.(*extensions.Ingress)
		ingEx := lbc.createIngress(ing)
		lbc.cnf.AddOrUpdateIngress(name, &ingEx)
	}
}

func (lbc *LoadBalancerController) enqueueIngressForService(obj interface{}) {
	svc := obj.(*api.Service)
	ings := lbc.getIngressesForService(svc)
	for _, ing := range ings {
		lbc.ingQueue.enqueue(&ing)
	}
}

func (lbc *LoadBalancerController) getIngressesForService(svc *api.Service) []extensions.Ingress {
	ings, err := lbc.ingLister.GetServiceIngress(svc)
	if err != nil {
		glog.V(3).Infof("ignoring service %v: %v", svc.Name, err)
		return nil
	}
	return ings
}

func (lbc *LoadBalancerController) getIngressForEndpoints(obj interface{}) []extensions.Ingress {
	var ings []extensions.Ingress
	endp := obj.(*api.Endpoints)
	svcKey := endp.GetNamespace() + "/" + endp.GetName()
	svcObj, svcExists, err := lbc.svcLister.Store.GetByKey(svcKey)
	if err != nil {
		glog.V(3).Infof("error getting service %v from the cache: %v\n", svcKey, err)
	} else {
		if svcExists {
			ings = append(ings, lbc.getIngressesForService(svcObj.(*api.Service))...)
		}
	}
	return ings
}

func (lbc *LoadBalancerController) createIngress(ing *extensions.Ingress) nginx.IngressEx {
	ingEx := nginx.IngressEx{
		Ingress: ing,
	}

	ingEx.Secrets = make(map[string]*api.Secret)
	for _, tls := range ing.Spec.TLS {
		secretName := tls.SecretName
		secret, err := lbc.client.Secrets(ing.Namespace).Get(secretName)
		if err != nil {
			glog.Warningf("Error retriveing secret %v for ing %v: %v", secretName, ing.Name, err)
			continue
		}
		ingEx.Secrets[secretName] = secret
	}

	ingEx.Endpoints = make(map[string]*api.Endpoints)
	if ing.Spec.Backend != nil {
		endps, err := lbc.getEndpointsForIngressBackend(ing.Spec.Backend, ing.Namespace)
		if err != nil {
			glog.V(3).Infof("Error retriviend endpoints for the services %v: %v", ing.Spec.Backend.ServiceName, err)
		} else {
			ingEx.Endpoints[ing.Spec.Backend.ServiceName] = endps
		}
	}

	for _, rule := range ing.Spec.Rules {
		if rule.IngressRuleValue.HTTP == nil {
			continue
		}

		for _, path := range rule.HTTP.Paths {
			endps, err := lbc.getEndpointsForIngressBackend(&path.Backend, ing.Namespace)
			if err != nil {
				glog.V(3).Infof("Error retriviend endpoints for the services %v: %v", path.Backend.ServiceName, err)
			} else {
				ingEx.Endpoints[path.Backend.ServiceName] = endps
			}

		}
	}

	return ingEx
}

func (lbc *LoadBalancerController) getEndpointsForIngressBackend(backend *extensions.IngressBackend, namespace string) (*api.Endpoints, error) {
	svcKey := namespace + "/" + backend.ServiceName
	svcObj, svcExists, err := lbc.svcLister.Store.GetByKey(svcKey)
	if err != nil {
		glog.V(3).Infof("error getting service %v from the cache: %v", svcKey, err)
		return nil, err
	}
	if svcExists {
		svc := svcObj.(*api.Service)
		endps, err := lbc.endpLister.GetServiceEndpoints(svc)
		if err != nil {
			glog.V(3).Infof("error getting endpoints for service %v from the cache: %v", svc, err)
			return nil, err
		}
		return &endps, nil
	}

	return nil, fmt.Errorf("Svc %s doesn't exists", svcKey)
}

func parseNGINXConfigMaps(nginxConfigMaps string) (string, string, error) {
	res := strings.Split(nginxConfigMaps, "/")
	if len(res) != 2 {
		return "", "", fmt.Errorf("NGINX configmaps name must follow the format <namespace>/<name>, got: %v", nginxConfigMaps)
	}
	return res[0], res[1], nil
}
