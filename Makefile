.PHONY: server build start

dev:
	go run src/*.go

install:
	go mod download

build-macos:
	make install && \
	CGO_ENABLED=0 GOOS=darwin go build -o ./build/auth-app ./src/*.go

build-linux:
	make install && \
	CGO_ENABLED=0 GOOS=linux go build -o ./build/auth-app ./src/*.go

start:
	./build/auth-app

deploy:
	docker buildx build --platform=linux/amd64 -t sifatulrabbi/finetrack-auth-app:latest -f Dockerfile --push .
