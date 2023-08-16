.PHONY: server build start

dev:
	go run src/*.go

install:
	go mod download

build-macos:
	make install && \
	GOOS=darwin go build -o ./build/auth-app ./src/*.go

build:
	make install && \
	GOOS=linux go build -o ./build/auth-app ./src/*.go

start:
	./build/auth-app
