FROM golang:1.20-alpine AS build
WORKDIR /go/src/service-log
COPY . .
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.6 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe
RUN CGO_ENABLED=0 go build -o /go/bin/service-log ./cmd/service-log

FROM scratch
COPY --from=build /go/bin/service-log /bin/service-log
COPY --from=build /bin/grpc_health_probe /bin/grpc_health_probe
ENTRYPOINT [ "/bin/service-log" ]