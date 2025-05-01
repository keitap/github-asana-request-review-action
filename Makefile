DOCKER_REPOSITORY = ghcr.io/keitap/github-asana-request-review-action:1.1.4

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: build
build: test lint
	docker buildx build --push --platform linux/amd64,linux/arm64 -t "${DOCKER_REPOSITORY}" .
