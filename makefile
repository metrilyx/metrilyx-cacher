
SHELL=/bin/bash

clean:
	rm -rf ./bin
	rm -rf ./build

build:
	go get -d -v ./...
	go build -v metrilyx-cacher.go 
	mkdir -p ./build/metrilyx-cacher/opt/metrilyx/bin
	mv metrilyx-cacher ./build/metrilyx-cacher/opt/metrilyx/bin/

all: clean build
