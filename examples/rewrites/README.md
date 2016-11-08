# Rewrites Support

To load balance an application that requires rewrites with NGINX Ingress controllers, you need to add the **nginx.org/rewrites** annotation to your Ingress resource definition. The annotation specifies which services need rewrites. The annotation syntax is as follows:
```
nginx.org/rewrites: "serviceName=service1 rewrite=/rewrite1/[;serviceName=service2 rewrite=/rewrite2/;...]"
```

In the following example we load balance two applications, which require rewrites:
```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: cafe-ingress
  annotations:
    nginx.org/rewrites: "serviceName=tea-svc rewrite=/;serviceName=coffee-svc rewrite=/beans/"
spec:
  rules:
  - host: cafe.example.com
    http:
      paths:
      - path: /tea/
        backend:
          serviceName: tea-svc
          servicePort: 80
      - path: /coffee/
        backend:
          serviceName: coffee-svc
          servicePort: 80
```

Requests to the tea service are rewritten as follows:

* /tea -> gets redirected to /tea/ first
* /tea/ -> /
* /tea/abc -> /abc

Requests to the coffee service are rewritten as follows:

* /coffee -> gets redirected to /coffee/ first
* /coffee/ -> /beans/
* /coffee/abc -> /beans/abc
