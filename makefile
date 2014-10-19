NAME = metrilyx-cacher
BUILD_PLATFORMS = linux/amd64

clean:
	rm -rvf ./bin

build-native:
	go build -o bin/$(NAME).$(shell go env GOOS).$(shell go env GOARCH)

build-platforms:
	for PLATFORM in $(BUILD_PLATFORMS); do \
		GO_OS=`echo $$PLATFORM | cut -d / -f 1`; \
		GO_ARCH=`echo $$PLATFORM | cut -d / -f 2`; \
		GOOS=$${GO_OS} GOARCH=$${GO_ARCH} go build -o bin/$(NAME).$${GO_OS}.$${GO_ARCH}; \
	done

build-all: clean build-native build-platforms

## Build for native platform only
build: build-native
