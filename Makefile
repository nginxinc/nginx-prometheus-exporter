VERSION = 0.8.0
TAG = $(VERSION)
PREFIX = nginx/nginx-prometheus-exporter

DOCKERFILEPATH = build
DOCKERFILE = Dockerfile

GIT_COMMIT = $(shell git rev-parse HEAD)

GOLANGCI_CONTAINER=golangci/golangci-lint:v1.29-alpine

export DOCKER_BUILDKIT = 1

.PHONY: nginx-prometheus-exporter
nginx-prometheus-exporter:
	GO111MODULE=on CGO_ENABLED=0 go build -mod=vendor -ldflags "-X main.version=$(VERSION) -X main.commit=$(GIT_COMMIT)" -o nginx-prometheus-exporter

.PHONY: lint
lint:
	docker run --rm \
	-v $(shell pwd):/go/src/github.com/nginxinc/nginx-prometheus-exporter \
	-w /go/src/github.com/nginxinc/nginx-prometheus-exporter \
	$(GOLANGCI_CONTAINER) golangci-lint run

.PHONY: test
test:
	GO111MODULE=on go test -mod=vendor ./...

.PHONY: container
container:
	docker build --build-arg VERSION=$(VERSION) --build-arg GIT_COMMIT=$(GIT_COMMIT) -f $(DOCKERFILEPATH)/$(DOCKERFILE) -t $(PREFIX):$(TAG) .

.PHONY: push
push: container
	docker push $(PREFIX):$(TAG)

.PHONY: clean
clean:
	-rm -r dist
	-rm nginx-prometheus-exporter
