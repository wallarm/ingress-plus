all: push

TAG = 0.1
PREFIX = gcr.io/nginx-ingress/nginx-ingress

nginx-ingress:
	CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo -ldflags '-w' -o nginx-ingress *.go

container: nginx-ingress
	docker build -t $(PREFIX):$(TAG) .

push: container
	gcloud docker push $(PREFIX):$(TAG)

osx:
	CGO_ENABLED=0 GOOS=darwin godep go build -a -installsuffix cgo -ldflags '-w' -o osx-nginx-ingress *.go

clean:
	rm -f nginx-ingress
