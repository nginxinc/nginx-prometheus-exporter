VERSION = 0.11.0
TAG = $(VERSION)
PREFIX = nginx/nginx-prometheus-exporter

.PHONY: nginx-prometheus-exporter
nginx-prometheus-exporter:
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X github.com/prometheus/common/version.Version=$(VERSION)" -o nginx-prometheus-exporter

.PHONY: build-goreleaser
build-goreleaser: ## Build all binaries using GoReleaser
	@goreleaser -v || (code=$$?; printf "\033[0;31mError\033[0m: there was a problem with GoReleaser. Follow the docs to install it https://goreleaser.com/install\n"; exit $$code)
	GOPATH=$(shell go env GOPATH) goreleaser build --clean --snapshot

.PHONY: lint
lint:
	docker run --pull always --rm -v $(shell pwd):/nginx-prometheus-exporter -w /nginx-prometheus-exporter -v $(shell go env GOCACHE):/cache/go -e GOCACHE=/cache/go -e GOLANGCI_LINT_CACHE=/cache/go -v $(shell go env GOPATH)/pkg:/go/pkg golangci/golangci-lint:latest golangci-lint --color always run

.PHONY: test
test:
	go test ./... -race -shuffle=on

.PHONY: container
container:
	docker build --build-arg VERSION=$(VERSION) --target container -f build/Dockerfile -t $(PREFIX):$(TAG) .

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
