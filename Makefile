MODULE = $(shell go list -m)
PACKAGES := $(shell go list ./... | grep -v /vendor/)

.PHONY: default
default: build

.PHONY: test
test: ## run unit tests
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg,$(PACKAGES), \
		go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg}; \
		tail -n +2 coverage.out >> coverage-all.out;)

.PHONY: test-cover
test-cover: test ## run unit tests and show test coverage information
	go tool cover -html=coverage-all.out

.PHONY: run
run: ## run the API server
	go run ${LDFLAGS} cmd/main/main.go

.PHONY: run-stop
run-stop: ## stop the API server
	@pkill -P `cat $(PID_FILE)` || true

.PHONY: build
build:  ## build the API server binary
	CGO_ENABLED=0 go build -a -o ./bin/crawler ./cmd/

.PHONY: build-docker
build-docker: ## build the API server as a docker image
	docker build -f build/Dockerfile -t crawler .

.PHONY: run-docker
run-docker: ## build the API server as a docker image
	docker run -it --rm -t crawler https://velikodny.com

.PHONY: clean
clean: ## remove temporary files
	rm -rf crawler coverage.out coverage-all.out
