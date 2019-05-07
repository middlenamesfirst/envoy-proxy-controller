FROM golang:1.12.4@sha256:83e8267be041b3ddf6a5792c7e464528408f75c446745642db08cfe4e8d58d18 as builder

WORKDIR /go/src/github.com/middlenamesfirst/envoy-proxy-controller

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /go/bin/envoy-proxy-controller .

FROM alpine:3.9.3@sha256:28ef97b8686a0b5399129e9b763d5b7e5ff03576aa5580d6f4182a49c5fe1913

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/bin/envoy-proxy-controller .

CMD ["./envoy-proxy-controller", "--kubeconfig=/root/config/kubeconfig.yaml"]
