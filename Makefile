VERSION = 0.9.0
TAG = $(VERSION)
PREFIX = nginx/nginx-prometheus-exporter

DOCKERFILEPATH = build
DOCKERFILE = Dockerfile

GIT_COMMIT = $(shell git rev-parse HEAD)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

export DOCKER_BUILDKIT = 1

.PHONY: nginx-prometheus-exporter
nginx-prometheus-exporter:
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(GIT_COMMIT) -X main.date=$(DATE)" -o nginx-prometheus-exporter

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run

.PHONY: test
test:
	go test ./...

.PHONY: container
container:
	docker build --build-arg VERSION=$(VERSION) --build-arg GIT_COMMIT=$(GIT_COMMIT) --build-arg DATE=$(DATE) --target container -f $(DOCKERFILEPATH)/$(DOCKERFILE) -t $(PREFIX):$(TAG) .

.PHONY: push
push: container
	docker push $(PREFIX):$(TAG)

.PHONY: deps
deps:
	@go mod tidy && go mod verify && go mod download

.PHONY: clean
clean:
	-rm -r dist
	-rm nginx-prometheus-exporter
