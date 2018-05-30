FROM golang:1.10 as builder
ARG VERSION
ARG GIT_COMMIT
WORKDIR /go/src/github.com/nginxinc/nginx-prometheus-exporter
COPY *.go ./
COPY vendor ./vendor
COPY collector ./collector
COPY client ./client
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT}" -o exporter .

FROM alpine:latest
COPY --from=builder /go/src/github.com/nginxinc/nginx-prometheus-exporter/exporter /usr/bin/
ENTRYPOINT [ "/usr/bin/exporter" ]