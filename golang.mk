.PHONY: deps
deps: tidy

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test:
	go test ./... -v -count=1 -coverprofile=profile.out -covermode=atomic ./... -timeout=5m

.PHONY: build
build:
	go build ./...