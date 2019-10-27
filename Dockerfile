FROM golang:1.13.3-alpine as build

COPY . /app

WORKDIR /app

RUN set -xe; \
    apk add git; \
    go build -o /go/bin/docker-proxy-go;

FROM alpine:3.10

ARG BUILD_DATE
ARG VCS_REF

LABEL maintainer="Philippe VANDERMOERE <philippe@wizacha.com" \
    org.label-schema.build-date=${BUILD_DATE} \
    org.label-schema.name="docker-proxy-go" \
    org.label-schema.vcs-ref=${VCS_REF} \
    org.label-schema.vcs-url="https://github.com/philippe-vandermoere/docker-proxy-go" \
    org.label-schema.schema-version="1.0.0"

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /go/bin/docker-proxy-go bin/docker-proxy-go

CMD bin/docker-proxy-go
