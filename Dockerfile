FROM docker.io/bitnami/minideb:buster
LABEL maintainer "Bitnami <containers@bitnami.com>"

ENV HOME="/" \
    OS_ARCH="amd64" \
    OS_FLAVOUR="debian-10" \
    OS_NAME="linux"

# Install required system packages and dependencies
COPY nginx-prometheus-exporter   /opt/bitnami/nginx-prometheus-exporter
RUN chmod g+rwX /opt/bitnami
RUN ln -sf /opt/bitnami/nginx-prometheus-exporter /usr/bin/exporter

ENV BITNAMI_APP_NAME="nginx-exporter" \
    BITNAMI_IMAGE_VERSION="0.9.0-debian-10-r140" \
    PATH="/opt/bitnami:$PATH"

EXPOSE 9113

WORKDIR /opt/bitnami
USER 1001
ENTRYPOINT [ "nginx-prometheus-exporter" ]
