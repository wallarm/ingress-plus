package main

import (
	"flag"
	"time"

	"github.com/golang/glog"

	"github.com/nginxinc/kubernetes-ingress/nginx-controller/controller"
	"github.com/nginxinc/kubernetes-ingress/nginx-controller/nginx"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

var (
	proxyURL = flag.String("proxy", "",
		`If specified, the controller assumes a kubctl proxy server is running on the
		given url and creates a proxy client. Regenerated NGINX configuration files
    are not written to the disk, instead they are printed to stdout. Also NGINX
    is not getting invoked. This flag is for testing.`)

	watchNamespace = flag.String("watch-namespace", api.NamespaceAll,
		`Namespace to watch for Ingress/Services/Endpoints. By default the controller
		watches acrosss all namespaces`)
)

func main() {
	flag.Parse()

	var kubeClient *client.Client
	var local = false

	if *proxyURL != "" {
		kubeClient = client.NewOrDie(&client.Config{
			Host: *proxyURL,
		})
		local = true
	} else {
		var err error
		kubeClient, err = client.NewInCluster()
		if err != nil {
			glog.Fatalf("Failed to create client: %v.", err)
		}
	}

	ngxc, _ := nginx.NewNGINXController("/etc/nginx/", local)
	ngxc.Start()
	lbc, _ := controller.NewLoadBalancerController(kubeClient, 30*time.Second, *watchNamespace, ngxc)
	lbc.Run()
}
