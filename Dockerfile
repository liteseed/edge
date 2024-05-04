FROM golang:latest

WORKDIR /app
COPY . .

RUN make docker
ENTRYPOINT ["./edge"]