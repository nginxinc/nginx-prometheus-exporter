VERSION = 1.3.0
TAG = $(VERSION)
PREFIX = nginx/nginx-prometheus-exporter
# renovate: datasource=docker depName=golangci/golangci-lint
GOLANGCI_LINT_VERSION = v1.60.3

.DEFAULT_GOAL:=nginx-prometheus-exporter

.PHONY: help
help: Makefile ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "; printf "Usage:\n\n    make \033[36m<target>\033[0m [VARIABLE=value...]\n\nTargets:\n\n"}; {printf "    \033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: nginx-prometheus-exporter
nginx-prometheus-exporter: ## Build nginx-prometheus-exporter binary
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X github.com/prometheus/common/version.Version=$(VERSION)" -o nginx-prometheus-exporter

.PHONY: build-goreleaser
build-goreleaser: ## Build all binaries using GoReleaser
	@goreleaser -v || (code=$$?; printf "\033[0;31mError\033[0m: there was a problem with GoReleaser. Follow the docs to install it https://goreleaser.com/install\n"; exit $$code)
	goreleaser build --clean --snapshot

.PHONY: lint
lint: ## Run linter
	docker run --pull always --rm -v $(shell pwd):/nginx-prometheus-exporter -w /nginx-prometheus-exporter -v $(shell go env GOCACHE):/cache/go -e GOCACHE=/cache/go -e GOLANGCI_LINT_CACHE=/cache/go -v $(shell go env GOPATH)/pkg:/go/pkg golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) golangci-lint --color always run

.PHONY: test
test: ## Run tests
	go test ./... -race -shuffle=on -coverprofile=coverage.txt -covermode=atomic

.PHONY: container
container: ## Build container image
	docker build --build-arg VERSION=$(VERSION) --target container -f build/Dockerfile -t $(PREFIX):$(TAG) .

.PHONY: push
push: container ## Push container image
	docker push $(PREFIX):$(TAG)

.PHONY: deps
deps: ## Add missing and remove unused modules, verify deps and download them to local cache
	@go mod tidy && go mod verify && go mod download

.PHONY: clean
clean: ## Clean up
	-rm -r dist
	-rm nginx-prometheus-exporter
