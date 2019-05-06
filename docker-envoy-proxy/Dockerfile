FROM envoyproxy/envoy-alpine:v1.10.0@sha256:930f6e7a1cde805d455683c9e54d452b5a6acc035ec77f86a7d499d6c92785a0

RUN mkdir -p /var/log/envoy

EXPOSE 9901

COPY envoy.yaml /etc/envoy/envoy.yaml
