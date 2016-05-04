.PHONY: clean build test clean

build:
	mkdir -p dist/
	go build -o dist/docker-machine-router

test:
	go list ./... | grep -v /vendor  | xargs go test

clean:
	rm -Rf dist/
