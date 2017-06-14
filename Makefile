.PHONY: build lint vet locally

default: build

build: vet gen compile

compile:
	go build -p 16 .

lint:
	golint .

vet:
	go vet .

gen:
	go generate

