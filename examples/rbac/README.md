# RBAC

For Kubernetes clusters with enabled [RBAC](https://kubernetes.io/docs/admin/authorization/rbac/), follow the steps below to deploy the Ingress controller:

1. Make sure you are a cluster admin.

1. If you would like to deploy the Ingress controller in a namespace other than `default`, change the namespace of the service account used in the cluster role binding in `nginx-ingress-rbac.yaml`. 

1. Create a service account, a cluster role and a cluster role binding for the Ingress controller:
	```
	$ kubectl create -f nginx-ingress-rbac.yaml
	```

1. As usual, create a secret with an SSL certificate and key for the default server of NGINX/NGINX Plus. It is recommended that you use your own certificate and key. 
    ```
    $ kubectl create -f default-server-secret.yaml
    ```

1. Deploy NGINX or NGINX Plus Ingress controller with the service account from the previous step:
	```
	$ kubectl create -f nginx-ingress-rc.yaml
	```
	or
	```
	$ kubectl create -f nginx-plus-ingress-rc.yaml
	```