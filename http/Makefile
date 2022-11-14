include ../golang.mk ../linting.mk

.PHONY: ci
ci: deps lint test

.PHONY: unit-test
unit-test:
	docker run --rm -v ${PWD}:/app -w /app andreykurennoy/golang:1.18.3-alpine3.15 \
	-c 'go test ./...'
