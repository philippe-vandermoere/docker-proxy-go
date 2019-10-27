FROM golang:1.13.3-alpine as build

COPY . /app

WORKDIR /app

RUN set -xe; \
    apk add git; \
    go build -o /go/bin/docker-proxy-go;

FROM alpine:3.10

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /go/bin/docker-proxy-go bin/docker-proxy-go

CMD bin/docker-proxy-go
