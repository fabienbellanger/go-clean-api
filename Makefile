.PHONY: all \
	install \
	update \
	update-all \
	format \
	lint \
	serve \
	serve-race \
	serve-logs \
	logs \
	watch \
	rabbit-client \
	rabbit-server \
	build \
	test \
	test-verbose \
	test-beautify \
	test-tparse \
	bench \
	clean \
	help \
	test-cover-count \
	cover-count \
	test-cover-atomic \
	cover-atomic \
	html-cover-count \
	html-cover-atomic \
	run-cover-count \
	run-cover-atomic \
	view-cover-count \
	view-cover-atomic

.DEFAULT_GOAL=help

include .env

# Read: https://kodfabrik.com/journal/a-good-makefile-for-go

# Go parameters
CURRENT_PATH=$(shell pwd)
MAIN_PATH=$(CURRENT_PATH)/cmd/main.go
GO_CMD=go
GO_INSTALL=$(GO_CMD) install
GO_RUN=$(GO_CMD) run
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
GO_MOD=$(GO_CMD) mod
GO_TOOL=$(GO_CMD) tool
GO_VET=$(GO_CMD) vet
GO_FMT=$(GO_CMD) fmt
BINARY_NAME=go-clean-api
BINARY_UNIX=$(BINARY_NAME)_unix

## all: Test and build application
all: test build

## install: Run go install
install:
	$(GO_INSTALL) ./...

## update: Update modules
update:
	$(GO_GET) -u ./... && $(GO_MOD) tidy

## update-all: Update all modules
update-all:
	$(GO_GET) -u ./... all && $(GO_MOD) tidy

## format: Run go fmt
format:
	$(GO_FMT) ./...

## lint: Run go vet
lint: format
	$(GO_VET) ./...

## serve: Serve API
serve:
	$(GO_RUN) $(MAIN_PATH) run

## serve-race: Serve API with -race option
serve-race:
	$(GO_RUN) run -race $(MAIN_PATH)

## serve-logs: Serve API with pretty logs
serve-logs:
	$(GO_RUN) $(MAIN_PATH) run | $(GO_RUN) $(MAIN_PATH) logs

## logs: Display server logs
logs:
	$(GO_RUN) $(MAIN_PATH) logs

## watch: Serve API with pretty logs and hot reload
watch:
	air | $(GO_RUN) $(MAIN_PATH) logs

## rabbit-client: Start RabbitMQ client
rabbit-client:
	$(GO_RUN) $(MAIN_PATH) rabbitmq -i client

## rabbit-server: Start RabbitMQ server
rabbit-server:
	$(GO_RUN) $(MAIN_PATH) rabbitmq -i server

build: format
	$(GO_BUILD) -ldflags "-s -w" -o $(BINARY_NAME) -v $(MAIN_PATH)

## test: Run test
test:
	$(GO_TEST) -cover ./...

## test-verbose: Run tests
test-verbose:
	$(GO_TEST) -cover -v ./...

## test-beautify: Run tests with gotestsum
test-beautify:
	gotestsum --format pkgname --debug

## test-tparse: Run tests with tparse
test-tparse:
	go test -cover -json ./... | tparse -trimpath -all

test-cover-count: 
	$(GO_TEST) -covermode=count -coverprofile=cover-count.out ./...

test-cover-atomic: 
	$(GO_TEST) -covermode=atomic -coverprofile=cover-atomic.out ./...

cover-count:
	$(GO_TOOL) cover -func=cover-count.out

cover-atomic:
	$(GO_TOOL) cover -func=cover-atomic.out

html-cover-count:
	$(GO_TOOL) cover -html=cover-count.out

html-cover-atomic:
	$(GO_TOOL) cover -html=cover-atomic.out

run-cover-count: test-cover-count cover-count
	rm cover-count.out
run-cover-atomic: test-cover-atomic cover-atomic
	rm cover-atomic.out
view-cover-count: test-cover-count html-cover-count
	rm cover-count.out
view-cover-atomic: test-cover-atomic html-cover-atomic
	rm cover-atomic.out

## bench: Run benchmarks
bench: 
	$(GO_TEST) -benchmem -bench=. ./...

## clean: Clean files
clean: 
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

help: Makefile
	@echo
	@echo "Choose a command run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo
