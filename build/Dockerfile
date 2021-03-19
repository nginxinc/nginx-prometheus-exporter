FROM golang:1.16 as builder
ARG VERSION
ARG GIT_COMMIT
WORKDIR /go/src/github.com/nginxinc/nginx-prometheus-exporter
COPY *.go ./
COPY vendor ./vendor
COPY go.mod go.sum ./
COPY collector ./collector
COPY client ./client
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -ldflags "-X main.version=${VERSION} -X main.commit=${GIT_COMMIT}" -o exporter .

FROM gcr.io/distroless/static-debian10:nonroot
COPY --from=builder /go/src/github.com/nginxinc/nginx-prometheus-exporter/exporter /usr/bin/

USER 1001:1001

ENTRYPOINT [ "/usr/bin/exporter" ]
