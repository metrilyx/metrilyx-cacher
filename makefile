
clean:
	rm -rvf ./bin
	rm -rvf ./build

.build:
	go install -v ./...
	[ -f ./build/opt/metrilyx/bin ] || mkdir -p ./build/opt/metrilyx/bin
	cp ../../../../bin/metrilyx-cacher ./build/opt/metrilyx/bin/
