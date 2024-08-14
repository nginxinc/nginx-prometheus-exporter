# syntax=docker/dockerfile:1.8
FROM golang:1.22 AS base
ARG VERSION
ARG TARGETARCH

WORKDIR /go/src/github.com/nginxinc/nginx-prometheus-exporter

FROM base AS builder
COPY --link go.mod go.sum ./
RUN go mod download
COPY --link *.go ./
COPY --link collector ./collector
COPY --link client ./client
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -trimpath -a -ldflags "-s -w -X main.version=${VERSION}" -o nginx-prometheus-exporter .


FROM scratch AS intermediate
COPY --from=base --link /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER 1001:1001
ENTRYPOINT [ "/usr/bin/nginx-prometheus-exporter" ]


FROM intermediate AS container
COPY --from=builder --link /go/src/github.com/nginxinc/nginx-prometheus-exporter/nginx-prometheus-exporter /usr/bin/


FROM intermediate AS goreleaser
ARG TARGETARCH
ARG TARGETVARIANT
ARG TARGETPLATFORM

LABEL org.nginx.exporter.image.build.target="${TARGETPLATFORM}"
LABEL org.nginx.exporter.image.build.version="goreleaser"

COPY --link dist/nginx-prometheus-exporter_linux_$TARGETARCH${TARGETVARIANT:+_7}*/nginx-prometheus-exporter /usr/bin/
