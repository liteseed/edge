all:
	go mod tidy
	go build -o ./dist/edge ./cmd
