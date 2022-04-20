.PHONY: build-docker
build-docker:
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t cubixle/groxy:latest --push .

.PHONY: build-docker-local
build-docker-local:
	docker buildx build --platform linux/amd64 --load -t cubixle/groxy:local .

