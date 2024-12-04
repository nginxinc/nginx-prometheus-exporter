# syntax=docker/dockerfile:1.12
FROM golang:1.23 AS builder
ARG VERSION
ARG TARGETARCH

WORKDIR /go/src/github.com/nginxinc/nginx-prometheus-exporter

COPY --link go.mod go.sum ./
RUN go mod download

COPY --link *.go ./
COPY --link collector ./collector
COPY --link client ./client

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -trimpath -a -ldflags "-s -w -X main.version=${VERSION}" -o nginx-prometheus-exporter .


FROM --platform=$BUILDPLATFORM alpine:3.20 AS certs

FROM scratch AS base
COPY --from=certs --link /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER 1001:1001
ENTRYPOINT [ "/usr/bin/nginx-prometheus-exporter" ]


FROM base AS container
COPY --from=builder --link /go/src/github.com/nginxinc/nginx-prometheus-exporter/nginx-prometheus-exporter /usr/bin/


FROM base AS goreleaser
ARG TARGETARCH
ARG TARGETVARIANT
ARG TARGETPLATFORM

LABEL org.nginx.exporter.image.build.target="${TARGETPLATFORM}"
LABEL org.nginx.exporter.image.build.version="goreleaser"

COPY --link dist/nginx-prometheus-exporter_linux_$TARGETARCH${TARGETVARIANT/v/_}*/nginx-prometheus-exporter /usr/bin/
