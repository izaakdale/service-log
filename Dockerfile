FROM golang:1.20-alpine AS build
WORKDIR /go/src/service-log
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/service-log ./cmd/service-log

FROM scratch
COPY --from=build /go/bin/service-log /bin/service-log
ENTRYPOINT [ "/bin/service-log" ]