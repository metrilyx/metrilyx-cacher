
SHELL=/bin/bash

clean:
	rm -rvf ./bin
	rm -rvf ./build

.build:
	go get -d -v ./...
	go install -v ./...
	mkdir -p ./build/metrilyx-cacher/opt/metrilyx/bin
	ls -la ../../../../
	cp ../../../../bin/metrilyx-cacher ./build/metrilyx-cacher/opt/metrilyx/bin/
