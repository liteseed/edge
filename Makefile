BINARY_NAME=edge

all:
	go mod tidy

build:
	GOARCH=amd64 GOOS=darwin go build -o ./dist/${BINARY_NAME}-darwin-amd64 ./cmd/main.go
	GOARCH=amd64 GOOS=linux go build -o ./dist/${BINARY_NAME}-linux-amd64 ./cmd/main.go
	GOARCH=386 GOOS=linux go build -o ./dist/${BINARY_NAME}-linux-386 ./cmd/main.go

run: build
	./dist/${BINARY_NAME}

clean:
	go clean
	rm ./dist/${BINARY_NAME}-darwin-amd64
	rm ./dist/${BINARY_NAME}-linux-amd64
	rm ./dist/${BINARY_NAME}-linux-386