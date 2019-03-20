VERSION = 0.3.0
PREFIX = nginx/nginx-prometheus-exporter
TAG = $(VERSION)
GIT_COMMIT = $(shell git rev-parse --short HEAD)

BUILD_DIR = build_output

nginx-prometheus-exporter: test
	CGO_ENABLED=0 go build -installsuffix cgo -ldflags "-X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT)" -o nginx-prometheus-exporter 

lint:
	golangci-lint run

test:
	go test ./...

container:
	docker build --build-arg VERSION=$(VERSION) --build-arg GIT_COMMIT=$(GIT_COMMIT) -t $(PREFIX):$(TAG) . 

push: container
	docker push $(PREFIX):$(TAG)

$(BUILD_DIR)/nginx-prometheus-exporter-linux-amd64:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -ldflags "-X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT)" -o $(BUILD_DIR)/nginx-prometheus-exporter-linux-amd64

$(BUILD_DIR)/nginx-prometheus-exporter-linux-i386:
	GOARCH=386 CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -ldflags "-X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT)" -o $(BUILD_DIR)/nginx-prometheus-exporter-linux-i386

release: $(BUILD_DIR)/nginx-prometheus-exporter-linux-amd64 $(BUILD_DIR)/nginx-prometheus-exporter-linux-i386
	mv $(BUILD_DIR)/nginx-prometheus-exporter-linux-amd64 $(BUILD_DIR)/nginx-prometheus-exporter && \
	tar czf $(BUILD_DIR)/nginx-prometheus-exporter-$(TAG)-linux-amd64.tar.gz -C $(BUILD_DIR) nginx-prometheus-exporter && \
	rm $(BUILD_DIR)/nginx-prometheus-exporter

	mv $(BUILD_DIR)/nginx-prometheus-exporter-linux-i386 $(BUILD_DIR)/nginx-prometheus-exporter && \
	tar czf $(BUILD_DIR)/nginx-prometheus-exporter-$(TAG)-linux-i386.tar.gz -C $(BUILD_DIR) nginx-prometheus-exporter && \
	rm $(BUILD_DIR)/nginx-prometheus-exporter
    
	shasum -a 256 $(BUILD_DIR)/nginx-prometheus-exporter-$(TAG)-linux-amd64.tar.gz $(BUILD_DIR)/nginx-prometheus-exporter-$(TAG)-linux-i386.tar.gz|sed "s|$(BUILD_DIR)/||" > $(BUILD_DIR)/sha256sums.txt

clean:
	-rm $(BUILD_DIR)/nginx-prometheus-exporter-$(TAG)-linux-amd64.tar.gz
	-rm $(BUILD_DIR)/nginx-prometheus-exporter-$(TAG)-linux-i386.tar.gz
	-rm $(BUILD_DIR)/sha256sums.txt
	-rmdir $(BUILD_DIR)
	-rm nginx-prometheus-exporter

