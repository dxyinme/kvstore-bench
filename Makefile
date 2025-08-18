.PHONY: all build clean test lint

all: build

build:
	@mkdir -p bin && cd cmd/kv-bench && go build -o ../../bin/kv-bench

clean:
	@rm -rf kv-bench

test:
	@go test -count 1 -v ./...

lint:
	@golangci-lint run

docker:
	@docker build -t kv-bench .