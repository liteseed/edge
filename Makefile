BINARY_NAME := edge
PKG := github.com/liteseed/edge
VERSION := 0.0.3

dev:
	go mod tidy
	go build -o ./build/dev/${BINARY_NAME} -ldflags="-X main.Version=${VERSION}-dev"  ./cmd/main.go

docker:
	go mod tidy
	go build -o ${BINARY_NAME} -ldflags="-X main.Version=${VERSION}" ./cmd/main.go

release:
	go mod tidy
	GOARCH=amd64 GOOS=darwin go build -o ./build/release/${BINARY_NAME}-darwin-amd64 -ldflags="-X main.Version=${VERSION}" ./cmd/main.go
	GOARCH=amd64 GOOS=linux go build -o ./build/release/${BINARY_NAME}-linux-amd64 -ldflags="-X main.Version=${VERSION}" ./cmd/main.go
	GOARCH=386 GOOS=linux go build -o ./build/release/${BINARY_NAME}-linux-386 -ldflags="-X main.Version=${VERSION}" ./cmd/main.go

clean:
	go clean
	rm -rf ./build
