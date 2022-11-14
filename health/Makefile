include ../golang.mk ../linting.mk

.PHONY: ci
ci: lint unit-test build

.PHONY: unit-test
unit-test:
	make test
