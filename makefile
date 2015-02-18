
SHELL=/bin/bash

clean:
	rm -rvf ./bin
	rm -rvf ./build

.build:
	go get -d -v ./...
	go install -v ./...
	ls -la ../../../../
	ls -la ../../../../bin
	mkdir -p ./build/metrilyx-cacher/opt/metrilyx/bin
	cp ../../../../bin/metrilyx-cacher ./build/metrilyx-cacher/opt/metrilyx/bin/
