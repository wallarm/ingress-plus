# Support for JSON Web Tokens (JWTs)

NGINX Plus supports validating JWTs with [ngx_http_auth_jwt_module](http://nginx.org/en/docs/http/ngx_http_auth_jwt_module.html). 

The Ingress controller provides the following 4 annotations for configuring JWT validation:

* Required: ```nginx.com/jwt-key: "secret"``` -- specifies a Secret resource with keys for validating JWTs. The keys must be stored in the `jwk` data field.
* Optional: ```nginx.com/jwt-realm: "realm"``` -- specifies a realm.
* Optional: ```nginx.com/jwt-token: "token"``` -- specifies a variable that contains JSON Web Token. By default, a JWT is expected in the `Authorization` header as a Bearer Token. 
* Optional: ```nginx.com/jwt-login-url: "url"``` -- specifies a URL to which a client is redirected in case of an invalid or missing JWT.

## Example

In the following example we enable JWT validation for the cafe-ingress Ingress:
```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: cafe-ingress
  annotations:
    nginx.com/jwt-key: "cafe-jwk" 
    nginx.com/jwt-realm: "Cafe App"  
    nginx.com/jwt-token: "$cookie_auth_token"
    nginx.com/jwt-login-url: "https://login.example.com"
spec:
  tls:
  - hosts:
    - cafe.example.com
    secretName: cafe-secret
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
* The keys must be deployed separately in the Secret `cafe-jwk`.
* The realm is  `Cafe App`.
* The token is extracted from the `auth_token` cookie.
* The login URL is `https://login.example.com`. 
