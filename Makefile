DOCKER_REPOSITORY = ghcr.io/keitap/github-asana-request-review-action:1.1.3

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	docker buildx build --push --platform linux/amd64,linux/arm64 -t "${DOCKER_REPOSITORY}" .
