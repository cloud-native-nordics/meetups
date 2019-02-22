export GO111MODULE=on
export GOOS:=$(shell go env GOOS)
export GOARCH:=$(shell go env GOARCH)

all: generate

generate: generator/bin/generator
	generator/bin/generator --config meetups.yaml --dry-run=false

dry-run: generator/bin/generator
	generator/bin/generator --config meetups.yaml --dry-run=true

validate: generator/bin/generator
	generator/bin/generator --config meetups.yaml --validate

build-docker:
	docker run -it -e GOOS=${GOOS} -e GOARCH=${GOARCH} -v $(shell pwd):/meetups -w /meetups golang:1.11 make bin-generator

generator/bin/generator bin-generator:
	cd generator && go build -mod vendor -o bin/generator .

clean:
	sudo rm generator/bin/generator
