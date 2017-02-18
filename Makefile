# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

PACKAGE = github.com/xackery/discordeq
SHORTNAME = discordeq
DOTPATH = github.com.xackery.discordeq

COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE = `date +%FT%T%z`
LDFLAGS = -ldflags "-X ${PACKAGE}/${SHORTNAME}lib.CommitHash=${COMMIT_HASH} -X ${PACKAGE}/${SHORTNAME}lib.BuildDate=${BUILD_DATE}"
NOGI_LDFLAGS = -ldflags "-X ${PACKAGE}/${SHORTNAME}lib.BuildDate=${BUILD_DATE}"

.PHONY: vendor docker check fmt lint test test-race vet test-cover-html help
.DEFAULT_GOAL := help

vendor: ## Install govendor and sync vendored dependencies
	go get github.com/kardianos/govendor
	govendor sync ${PACKAGE}

build: vendor ## Build binary
	go build ${LDFLAGS} ${PACKAGE}

build-race: vendor ## Build smite binary with race detector enabled
	go build -race ${LDFLAGS} ${PACKAGE}

install: vendor ## Install smite binary
	go install ${LDFLAGS} ${PACKAGE}

build-no-gitinfo: LDFLAGS = ${NOGI_LDFLAGS}
build-no-gitinfo: vendor ${SHORTNAME} ## Build without git info

docker: ## Build Docker container
	docker build -t ${SHORTNAME} .
	docker rm -f ${SHORTNAME}-build || true
	docker run --name ${SHORTNAME}-build ${SHORTNAME} ls /go/bin
	docker cp ${SHORTNAME}-build:/go/bin/${SHORTNAME} .
	docker rm ${SHORTNAME}-build

check: test-race test386 fmt vet ## Run tests and linters

test386: ## Run tests in 32-bit mode
	GOARCH=386 govendor test +local

test: ## Run tests
	govendor test +local

test-race: ## Run tests with race detector
	govendor test -race +local

fmt: ## Run gofmt linter
	@for d in `govendor list -no-status +local | sed 's/${DOTPATH}/./'` ; do \
		if [ "`gofmt -l $$d/*.go | tee /dev/stderr`" ]; then \
			echo "^ improperly formatted go files" && echo && exit 1; \
		fi \
	done

lint: ## Run golint linter
	@for d in `govendor list -no-status +local | sed 's/${DOTPATH}/./'` ; do \
		if [ "`golint $$d | tee /dev/stderr`" ]; then \
			echo "^ golint errors!" && echo && exit 1; \
		fi \
	done

vet: ## Run go vet linter
	@if [ "`govendor vet +local | tee /dev/stderr`" ]; then \
		echo "^ go vet errors!" && echo && exit 1; \
	fi

test-cover-html: PACKAGES = $(shell govendor list -no-status +local | sed 's/${DOTPATH}/./')
test-cover-html: ## Generate test coverage report
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		govendor test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
