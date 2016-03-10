# Examples

## Prerequisites

* Kubernetes 1.2 (TSL support for Ingress has been added in 1.2)
* For NGINX Plus: you've built and made available in your cluster
[NGINX Plus](https://github.com/nginxinc/kubernetes-ingress/tree/master/nginx-plus-controller) Controller
image and updated the container image field in the ```nginx-plus-ingress-rc.yaml``` file accordingly.

## Running the examples

1. Create coffee and tea services and replication controllers:

  ```
  $ kubectl create -f tea-rc.yaml
  $ kubectl create -f tea-svc.yaml
  $ kubectl create -f coffee-rc.yaml
  $ kubectl create -f coffee-svc.yaml
  ```
1. Create a Secret with an SSL certificate and a key:
  ```
  $ kubectl create -f cafe-secret.yaml
  ```

1. Create Ingress Resource:
  ```
  $ kubectl create -f cafe-ingress.yaml
  ```

1. Create either NGINX or NGINX Plus Ingress Controller:
  ```
  $ kubectl create -f nginx-ingress-rc.yaml
  ```
  or
  ```
  $ kubectl create -f nginx-plus-ingress-rc.yaml
  ```
  If you're creating the Plus controller, please make sure that the Docker image
  is available in your cluster.

1. The Controller container exposes ports 80, 443 (and 8080 for NGINX Plus )
on the host it runs. Make sure to add a firewall to allow incoming traffic
on this ports.

1. Find out the external IP address of the node of the controller:
  ```
  $ kubectl get pods -o wide
  NAME                          READY     STATUS    RESTARTS   AGE       NODE
  coffee-rc-mtjuw               1/1       Running   0          3m        kubernetes-minion-iikt
  coffee-rc-mu9ns               1/1       Running   0          3m        kubernetes-minion-cm0y
  nginx-plus-ingress-rc-86kkq   1/1       Running   0          1m        kubernetes-minion-iikt
  tea-rc-7w3fq                  1/1       Running   0          3m        kubernetes-minion-iikt
  ```

  ```
  $ kubectl get node kubernetes-minion-iikt -o json | grep -A 2 ExternalIP
      "type": "ExternalIP",
    "address": "XXX.YYY.ZZZ.III"
    }
  ```


1. We'll use ```curl```'s --insecure to turn off certificate verification of our self-signed
certificate and --resolve option to set the Host header of a request with ```cafe.example.com```
  To get coffee:
  ```
  $ curl --resolve cafe.example.com:443:XXX.YYY.ZZZ.III https://cafe.example.com/coffee --insecure
  <!DOCTYPE html>
  <html>
  <head>
  <title>Hello from NGINX!</title>
  <style>
      body {
          width: 35em;
          margin: 0 auto;
          font-family: Tahoma, Verdana, Arial, sans-serif;
      }
  </style>
  </head>
  <body>
  <h1>Hello!</h1>
  <h2>URI = /coffee</h2>
  <h2>My hostname is coffee-rc-mu9ns</h2>
  <h2>My address is 10.244.0.3:80</h2>
  </body>
  </html>
  ```
  To get tea:
  ```
  $ curl --resolve cafe.example.com:443:XXX.YYY.ZZZ.III https://cafe.example.com/tea --insecure
  <!DOCTYPE html>
  <html>
  <head>
  <title>Hello from NGINX!</title>
  <style>
      body {
          width: 35em;
          margin: 0 auto;
          font-family: Tahoma, Verdana, Arial, sans-serif;
      }
  </style>
  </head>
  <body>
  <h1>Hello!</h1>
  <h2>URI = /tea</h2>
  <h2>My hostname is tea-rc-w7rjr</h2>
  <h2>My address is 10.244.0.5:80</h2>
  </body>
  </html>
  ```
