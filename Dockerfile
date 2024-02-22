FROM golang:latest


WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o ./build/bungo ./cmd

EXPOSE 8080

CMD [ "./build/bungo" ]

