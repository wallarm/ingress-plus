# Deploying NGINX and NGINX Plus Controllers as a Daemon Set

You can deploy the NGINX or NGINX Plus controller as a [Daemon Set](http://kubernetes.io/docs/admin/daemons/). This allows you to deploy the controller on all or select nodes of your cluster.

To deploy the NGINX controller, run:
```
  $ kubectl create -f nginx-ingress.yaml
```

To deploy the NGINX Plus controller, run:
```
  $ kubectl create -f nginx-plus-ingress.yaml
```

Once deployed, by default, a controller pod is running on every node of the cluster. The pods are accessible through ports 80 and 443 of each node they get scheduled on.

Optionally, you can choose to run the controller pods on only select nodes. To accomplish this:
1. Add a label to each node on which you want to run a controller pod. For example:
  ```
  kubectl label node node-1 role=nginx-ingress
  kubectl label node node-2 role=nginx-ingress
  ```
  where *node-1* and *node-2* are some nodes of your cluster.

1. Uncomment the **nodeSelector** related lines (11-12) in the corresponding daemon set yaml file and specify a label to use to select nodes (`role=nginx-ingress` in this example).

1. Deploy the controller. The pods are scheduled only on *node-1* and *node-2*.
