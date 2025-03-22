BINARY_NAME=pumpfun-monitor
BUILD_FOLDER=./build
SRC_FOLDER=./src
GO_FILES=$(wildcard *.go)
PREFIX=/usr/local
DESTDIR=

.PHONY: build run clean test mod-tidy help install

build:
	@mkdir -p ${BUILD_FOLDER}
	go build -o ${BUILD_FOLDER}/${BINARY_NAME} ${SRC_FOLDER}/${GO_FILES}

run: build
	@./${BUILD_FOLDER}/${BINARY_NAME} $(ARGS)

clean:
	@go clean
	@rm -rf ${BUILD_FOLDER}

test:
	@go test -v ./...

mod-tidy:
	@go mod tidy

install: build
	@mkdir -p ${DESTDIR}${PREFIX}/bin
	install -m 755 ${BUILD_FOLDER}/${BINARY_NAME} ${DESTDIR}${PREFIX}/bin/${BINARY_NAME}


help:
	@echo "Targets:"
	@echo "  build    - Compile project (all .go files)"
	@echo "  run      - Build and execute"
	@echo "  clean    - Remove build artifacts"
	@echo "  test     - Run tests verbosely"
	@echo "  mod-tidy - Clean up dependencies"
	@echo "  install  - Install binary to system location (default: /usr/local/bin)"
