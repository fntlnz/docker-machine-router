.PHONY: clean build

build:
	mkdir -p dist/
	go build -o dist/docker-machine-router
