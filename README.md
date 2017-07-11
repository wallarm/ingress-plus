# NGINX Ingress Controller

This repo provides an implementation of an Ingress controller for NGINX and NGINX Plus. This implementation is different from the NGINX Ingress controller in [kubernetes/ingress repo](https://github.com/kubernetes/ingress).

## What is Ingress?

An Ingress is a Kubernetes resource that lets you configure an HTTP load balancer for your Kubernetes services. Such a load balancer usually exposes your services to clients outside of your Kubernetes cluster. An Ingress resource supports:
* Exposing services:
    * Via custom URLs (for example, service A at the URL `/serviceA` and service B at the URL `/serviceB`).
    * Via multiple host names (for example, `foo.example.com` for one group of services and `bar.example.com` for another group).
* Configuring SSL termination for each exposed host name.

See the [Ingress User Guide](http://kubernetes.io/docs/user-guide/ingress/) to learn more.

## What is an Ingress Controller?

An Ingress controller is an application that monitors Ingress resources via the Kubernetes API and updates the configuration of a load balancer in case of any changes. Different load balancers require different Ingress controller implementations. Typically, an Ingress controller is deployed as a pod in a cluster. In the case of software load balancers, such as NGINX, an Ingress controller is deployed in a pod along with a load balancer.

See https://github.com/kubernetes/contrib/tree/master/ingress/controllers/ to learn more about Ingress controllers and find out about different implementations.

## NGINX Ingress Controller

We provide an Ingress controller for NGINX and NGINX Plus that supports the following Ingress features:
* SSL termination
* Path-based rules
* Multiple host names

We provide the following extensions to our Ingress controller:
* [Websocket](examples/websocket), which allows you to load balance Websocket applications.
* [SSL Services](examples/ssl-services), which allows you to load balance HTTPS applications.
* [Rewrites](examples/rewrites), which allows you to rewrite the URI of a request before sending it to the application.
* [Session Persistence](examples/session-persistence) (NGINX Plus only), which guarantees that all the requests from the same client are always passed to the same backend container.

Additionally, we provide a mechanism to customize the NGINX configuration. Refer to the [examples folder](examples) to find out how to [deploy](examples/complete-example) the Ingress controller and [customize](examples/customization) the NGINX configuration.

## Benefits of Using the Ingress Controller with NGINX Plus

[NGINX Plus](https://www.nginx.com/products/) is a commercial version of NGINX that comes with advanced features and support.

The Ingress controller leverages the advanced features of NGINX Plus, which gives you the following additional benefits:

* **Reduced number of configuration reloads**
Every time the number of pods of services you expose via Ingress changes, the Ingress controller updates the configuration of NGINX to reflect those changes. For the open source NGINX software, the configuration file must be changed and the configuration reloaded. For NGINX Plus, the [on-the-fly reconfiguration](https://www.nginx.com/products/on-the-fly-reconfiguration/) feature is utilized, which allows NGINX Plus to be updated on-the-fly without reloading the configuration. This prevents a potential increase of memory usage and overall system overloading, which could occur with too frequent configuration reloads.
* **Real-time statistics**
NGINX Plus provides you with [advanced statistics](https://www.nginx.com/products/live-activity-monitoring/), which you can access either through the API or via the built-in dashboard. This can give you insights into how NGINX Plus and your applications are performing.
* **Session persistence** When enabled, NGINX Plus makes sure that all the requests from the same client are always passed to the same backend container using the *sticky cookie* method. Refer to the [session persistence examples](examples/session-persistence) to find out how to configure it.

**Note**: Deployment of the Ingress controller for NGINX Plus requires you to do one extra step: build your own [Docker image](nginx-controller) using the certificate and key for your subscription.
The Docker image of the Ingress controller for NGINX is [available on Docker Hub](https://hub.docker.com/r/nginxdemos/nginx-ingress/).

## Using Multiple Ingress Controllers

You can run multiple Ingress controllers at the same time. For example, if your Kubernetes cluster is deployed in cloud, you can run the NGINX controller and the corresponding cloud HTTP load balancing controller. Refer to the [example](examples/multiple-ingress-controllers) to learn more.

## Advanced Load Balancing (Beyond Ingress)

When your requirements go beyond what Ingress offers, you can use NGINX and
NGINX Plus without the Ingress Controller.

NGINX Plus comes with a [DNS-based dynamic reconfiguration feature](https://www.nginx.com/blog/dns-service-discovery-nginx-plus/), which lets you keep the list of the endpoints of your services in sync with NGINX Plus. Read more about how to setup NGINX Plus this way in [Load Balancing Kubernetes Services with NGINX Plus](https://www.nginx.com/blog/load-balancing-kubernetes-services-nginx-plus/).

## Production Status

This is the preview version of the Ingress controller.

## Support

Support from the [NGINX Professional Services Team](https://www.nginx.com/services/) is available when using the NGINX Plus Ingress controller.

## Contacts

We’d like to hear your feedback! If you have any suggestions or experience issues with our Ingress controller, please create an issue or send a pull request on Github.
You can contact us directly via [kubernetes@nginx.com](mailto:kubernetes@nginx.com).
