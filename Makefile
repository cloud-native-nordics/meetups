

all: generate

generate: generator/bin/generator
	generator/bin/generator --config meetups.yaml --dry-run=false

validate: generator/bin/generator
	generator/bin/generator --config meetups.yaml --validate

generator/bin/generator bin-generator:
	docker run -it -v $(shell pwd)/generator:/generator -w /generator golang:1.11 /bin/bash -c "go build -mod vendor -o bin/generator ."

clean:
	sudo rm generator/bin/generator
