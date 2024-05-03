FROM golang:latest

WORKDIR /app
COPY config.json go.mod ./

RUN make docker

CMD ["./edge", "start"]
EXPOSE 8080