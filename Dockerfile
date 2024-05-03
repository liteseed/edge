FROM golang:latest

WORKDIR /app
COPY . .

RUN make docker
RUN ./edge generate
CMD ["./edge", "start"]