# Customization of NGINX Configuration

You can customize the NGINX configuration using ConfigMaps or Annotations. For now, you can set the values of the following
NGINX directives:
* [proxy_connect_timeout](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_connect_timeout)
* [proxy_read_timeout](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_read_timeout)
* [client_max_body_size](http://nginx.org/en/docs/http/ngx_http_core_module.html#client_max_body_size)
* [server_names_hash_max_size](http://nginx.org/en/docs/http/ngx_http_core_module.html#server_names_hash_max_size) via ConfigMaps only
* [server_names_hash_bucket_size](http://nginx.org/en/docs/http/ngx_http_core_module.html#server_names_hash_bucket_size) via ConfigMaps only

## Using ConfigMaps

1. Make sure that you specify the configmaps resource to use when you start an Ingress Controller.
For example, `-nginx-configmaps=default/nginx-config`, where we specify
the config map to use with the following format: `<namespace>/<name>`. See [nginx-ingress-rc.yaml](../complete-example/nginx-ingress-rc.yaml) or
[nginx-plus-ingress-rc.yaml](../complete-example/nginx-plus-ingress-rc.yaml) files.

1. Create a configmaps file with the name *nginx-config.yaml* and set the values
that make sense for your setup:
  ```yaml
  kind: ConfigMap
  apiVersion: v1
  metadata:
    name: nginx-config
  data:
    proxy-connect-timeout: "10s"
    proxy-read-timeout: "10s"
    client-max-body-size: "2m"
  ```
  See the **nginx-config.yaml** from this directory for a complete example.

2. Create a configmaps resource:
  ```
  $ kubectl create -f nginx-config.yaml
  ```
  The NGINX configuration will be updated.

3. If you want to update the configmaps, update the file and replace the config map:
  ```
  $ kubectl replace -f nginx-config.yaml
  ```
  The NGINX configuration will be updated.

## Using Annotations

If you want to customize the configuration for a particular Ingress resource only, you can use Annotations.
Here is an example (**cafe-ingress-with-annotations.yaml**):
```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: cafe-ingress-with-annotations
  annotations:
    nginx.org/proxy-connect-timeout: "30s"
    nginx.org/proxy-read-timeout: "20s"
    nginx.org/client-max-body-size: "4m"
spec:
  rules:
  - host: cafe.example.com
    http:
      paths:
      - path: /tea
        backend:
          serviceName: tea-svc
          servicePort: 80
      - path: /coffee
        backend:
          serviceName: coffee-svc
          servicePort: 80
```
Annotations take precedence over ConfigMaps.
