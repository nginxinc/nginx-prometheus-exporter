FROM golang:1.16-alpine as builder
ARG VERSION
ARG GIT_COMMIT
ARG TARGETARCH

WORKDIR /go/src/github.com/nginxinc/nginx-prometheus-exporter
COPY *.go ./
COPY vendor ./vendor
COPY go.mod go.sum ./
COPY collector ./collector
COPY client ./client
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -mod=vendor -a -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${GIT_COMMIT}" -o nginx-prometheus-exporter .


FROM gcr.io/distroless/static-debian10:nonroot as base
USER 1001:1001
ENTRYPOINT [ "/usr/bin/nginx-prometheus-exporter" ]


FROM base as container
COPY --from=builder /go/src/github.com/nginxinc/nginx-prometheus-exporter/nginx-prometheus-exporter /usr/bin/


FROM base as local
ARG TARGETARCH
ARG TARGETVARIANT

COPY dist/nginx-prometheus-exporter_linux_$TARGETARCH${TARGETVARIANT:+_7}/nginx-prometheus-exporter /usr/bin
