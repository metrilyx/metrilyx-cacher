
SHELL=/bin/bash

clean:
	rm -rvf ./bin
	rm -rvf ./build

.build:
	go get -d -v ./...
	go build -v metrilyx-cacher.go 
	mkdir -p ./build/metrilyx-cacher/opt/metrilyx/bin
	cp metrilyx-cacher ./build/metrilyx-cacher/opt/metrilyx/bin/
