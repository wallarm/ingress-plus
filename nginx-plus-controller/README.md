# NGINX Plus Ingress Controller

This is an implementation of a Kubernetes Ingress controller for NGINX Plus, which provides HTTP load balancing for applications your deploy in your Kubernetes cluster. You can find more details on what an Ingress controller is as well as difference between NGINX and NGINX Plus controllers on the [main page](https://github.com/nginxinc/kubernetes-ingress).

## How to Use the Controller

To find examples on how to deploy, configure and use the NGINX Plus Ingress controller, please see [the examples folder](../examples). The examples require the Docker image of the controller to be available to your Kubernetes cluster. You need to build the image. Please read the next section for instructions on how to do it.

### Prerequisites

Before you can build the image, make sure that the following software is installed on your machine:
* [Docker](https://www.docker.com/products/docker)
* [GNU Make](https://www.gnu.org/software/make/)

Additionally, you must have the NGINX Plus license. If you don't have one, you can sign up for a [free 30-day trial](https://www.nginx.com/free-trial-request/).  Put the certificate (`nginx-repo.crt`) and the key (`nginx-repo.key`) of your license inside this folder.

### Building the image

We build the image using the make utility. The **Makefile** we provide has the following targets:
* **test**: runs unit tests.
* **nginx-plus-ingress**: creates the controller binary.
* **container**: builds a Docker image.
* **push**: pushes the image to the private Docker registry.
* **all** (the default target): executes the four targets above in the order listed. If one of the targets fails, the execution process stops, reporting an error.

The **Makefile** contains the following main variables, which you should customize (either by changing the Makefile or by overriding the variables in the make command):
* **PREFIX** -- the name of the image. For example, `nginx-plus-ingress`
* **TAG** -- the tag added to the image. For example, `0.3`
* **PUSH_TO_GCR** . If you’re running your Kubernetes in GCE and using Google Container Registry, make sure that `PUSH_TO_GCR = 1`. This means using the `gcloud docker push` command to push the image, which is convenient when pushing images to GCR. By default, the variable is unset and the regular `docker push` command is used to push the image to the registry.

Let’s create the controller binary, build an image and push the image to the private registry. Make sure to run the `docker login` command first to login to the registry. If you’re using Google Container Registry, as we are in our example here, you don’t need to use the docker command to login. However, make sure you’re logged into the gcloud tool (using the `gcloud auth login` command). 

In this folder we run the following commands in the shell:
```
$ make clean
$ make PREFIX=gcr.io/my-k8s-project/nginx-plus-ingress TAG=0.3 PUSH_TO_GCR=1
```

Where **my-k8s-project** is the name of the GCE project where we run our Kubernetes cluster. As the result, the image -- **gcr.io/my-k8s-project/nginx-plus-ingress:0.3** --  is built and pushed to the registry.

By default, to compile the controller we use the [golang](https://hub.docker.com/_/golang/) container that we run as part of the building process. If you want to compile the controller using your local golang environment, specify `BUILD_IN_CONTAINER=0` when you run the make command.
