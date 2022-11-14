PACKAGES := `find . -name "go.mod" -exec dirname {} \;`
pwd:=$(PWD)
coverage_path="$(pwd)/coverage.out"

.PHONY: all-ci
all-ci:
	@for i in ${PACKAGES}; do \
		cd $$i && make ci || exit 1 && cd $(pwd); \
		cat $$i/profile.out >> ${coverage_path}; \
	done

.PHONY: all-build
all-build:
	@for i in ${PACKAGES}; do \
		cd $$i && make build || exit 1 && cd $(pwd); \
	done

.PHONY: all-deps
all-deps:
	@for i in ${PACKAGES}; do \
		cd $$i && make deps || exit 1 && cd $(pwd); \
	done

.PHONY: all-lint
all-lint:
	@for i in ${PACKAGES}; do \
		cd $$i && make lint || exit 1 && cd $(pwd); \
	done

.PHONY: all-test
all-test:
	@for i in ${PACKAGES}; do \
		cd $$i && make test || exit 1 && cd $(pwd); \
	done
