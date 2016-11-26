# Changelog

### 0.6.0

* [75](https://github.com/nginxinc/kubernetes-ingress/pull/75): Add the HSTS settings in the configmap and annotations. Thanks to [Nico Schieder](https://github.com/thetechnick).
* [74](https://github.com/nginxinc/kubernetes-ingress/pull/74): Fix the issue of the `kubernetes.io/ingress.class` annotation handling. Thanks to [Tang Le](https://github.com/tangle329).
* [70](https://github.com/nginxinc/kubernetes-ingress/pull/70): Add support for the alpine-based image for the NGINX controller.
* [68](https://github.com/nginxinc/kubernetes-ingress/pull/68): Support for proxy-buffering settings in the configmap and annotations. Thanks to [Mark Daniel Reidel](https://github.com/df-mreidel).
* [66](https://github.com/nginxinc/kubernetes-ingress/pull/66): Support for custom log-format in the configmap. Thanks to [Mark Daniel Reidel](https://github.com/df-mreidel).
* [65](https://github.com/nginxinc/kubernetes-ingress/pull/65): Add HTTP/2 as an option in the configmap and annotations. Thanks to [Nico Schieder](https://github.com/thetechnick).
* The NGINX Plus controller image is now based on Ubuntu Xenial.

### 0.5.0

* Update NGINX version to 1.11.5.
* [64](https://github.com/nginxinc/kubernetes-ingress/pull/64): Add the `nginx.org/rewrites` annotation, which allows to rewrite the URI of a request before sending it to the application. Thanks to [Julian](https://github.com/jmastr).
* [62](https://github.com/nginxinc/kubernetes-ingress/pull/62): Add the `nginx.org/ssl-services` annotation, which allows load balancing of HTTPS applications. Thanks to [Julian](https://github.com/jmastr).

### 0.4.0

* [54](https://github.com/nginxinc/kubernetes-ingress/pull/54): Previously, when specifying the port of a service in an Ingress rule, you had to use the value of the target port of that port of the service, which was incorrect. Now you must use the port value or the name of the port of the service instead of the target port value. **Note**: Please make necessary changes to your Ingress resources, if ports of your services have different values of the port and the target port fields.
* [55](https://github.com/nginxinc/kubernetes-ingress/pull/55): Add support for the `kubernetes.io/ingress.class` annotation in Ingress resources.
* [58](https://github.com/nginxinc/kubernetes-ingress/pull/58): Add the version information to the controller. For each version of the NGINX controller, you can find a corresponding image on [DockerHub](https://hub.docker.com/r/nginxdemos/nginx-ingress/tags/) with a tag equal to the version. The latest version is available through the `latest` tag.

The previous version was 0.3


### Notes

* Except when mentioned otherwise, the controller refers both to the NGINX and the NGINX Plus Ingress controllers.
