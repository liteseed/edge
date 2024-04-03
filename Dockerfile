FROM golang:latest

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o ./build/edge ./cmd

EXPOSE 8080

CMD [ "./build/edge" ]
