# NGINX Plus Ingress Controller

## Building the Image

1. You must have the following software installed on your machine to build the image:
	* [go lang](https://golang.org/dl/)
	* [Docker machine](https://docs.docker.com/machine/)

	You must also have an NGINX Plus license. If you don't have one, you can sign up for a [free 30-day trial](https://www.nginx.com/#free-trial).

	Put the key and the certificate file of your license inside this folder.

1. Change the ```Makefile``` according to your environment.

	If you're running Kubernetes in GCE, change the ```PREFIX``` variable  to ```PREFIX = gcr.io/<project-id>/nginx-plus-ingress```, where ```project-id```is the ID of your GCE project.


1. Run the ```make``` command to build the image and upload it to the private Docker registry you use:
	```
	$ make
	```

## Requirements for your Kubernetes setup

[DNS cluster addon](https://github.com/kubernetes/kubernetes/tree/master/cluster/addons/dns) must be enabled. It is enabled by default in deployments for most cloud providers.
NGINX Plus gets the IP address of a service or IP addresses of Endpoints for a headless service by resolving
the service DNS name.
