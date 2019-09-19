# Docker-proxy go

![CircleCI](https://img.shields.io/circleci/build/github/philippe-vandermoere/docker-proxy-go)
[![codecov](https://codecov.io/gh/philippe-vandermoere/docker-proxy-go/branch/master/graph/badge.svg)](https://codecov.io/gh/philippe-vandermoere/docker-proxy-go)

Go implementation of [docker-proxy](https://github.com/philippe-vandermoere/docker-proxy).

## Development

### Start

Run project

```bash
docker-compose run golang go run main.go
```

### Lint

```bash
docker-compose run golang bin/lint
```

### Test

```bash
docker-compose run golang bin/test
```
