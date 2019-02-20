GO111MODULE=on
export GO111MODULE

all: generate

generate: generator/bin/generator
	generator/bin/generator --config meetups.yaml --dry-run=false

dry-run: generator/bin/generator
	generator/bin/generator --config meetups.yaml --dry-run=true

validate: generator/bin/generator
	generator/bin/generator --config meetups.yaml --validate

build-docker:
	docker run -it -v $(shell pwd):/meetups -w /meetups golang:1.11 make bin-generator

generator/bin/generator bin-generator:
	cd generator && go build -mod vendor -o bin/generator .

clean:
	sudo rm generator/bin/generator
