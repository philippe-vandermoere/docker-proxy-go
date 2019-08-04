FROM golang:1.12.8-alpine as build

COPY . /app

WORKDIR /app

RUN set -xe; \
    apk add git; \
    go build -o /go/bin/docker-proxy-go;

FROM golang:1.12.8-alpine

WORKDIR /app

COPY --from=build /go/bin/docker-proxy-go bin/docker-proxy-go
COPY ./template ./template

CMD bin/docker-proxy-go
