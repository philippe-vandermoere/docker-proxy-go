default:

test:
    bin/test

lint:
	bin/lint

docker_build:
	docker build . \
		--build-arg VCS_REF=$(git rev-parse --short HEAD) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
