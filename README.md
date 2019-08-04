# docker-proxy-go

Go implementation of [docker-proxy](https://github.com/philippe-vandermoere/docker-proxy).

## Development

### Start

Run project

```bash
docker-compose run golang go run main.go
```

### Lint

```bash
docker-compose run golang gofmt -d -e -s .
```

### Test

```bash
docker-compose run golang go test
```
