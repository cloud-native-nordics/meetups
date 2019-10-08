export GO111MODULE=on
export GOOS:=$(shell go env GOOS)
export GOARCH:=$(shell go env GOARCH)

all: generate

generate: bin/generator
	bin/generator --dry-run=false

dry-run: bin/generator
	bin/generator --dry-run=true

validate: bin/generator
	bin/generator --validate

stats: bin/generator
	bin/generator --stats

test:
	cd generator && go test

build-docker:
	docker run -it -e GOOS=${GOOS} -e GOARCH=${GOARCH} -v $(shell pwd):/meetups -w /meetups golang:1.12 make bin-generator

bin/generator bin-generator:
	go build -mod vendor -o bin/generator $(shell pwd)/generator/...

clean:
	sudo rm bin/generator

pre-commit:
	$(MAKE) build-docker
	$(MAKE) generate
	$(MAKE) validate
	$(MAKE) test
	gofmt -s -w generator/*.go
