FROM alpine:latest

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

WORKDIR /bungo

VOLUME ["/bungo/data"]

COPY cmd/bungo /bungo/bungo
EXPOSE 8080

ENTRYPOINT [ "/bungo/bungo" ]