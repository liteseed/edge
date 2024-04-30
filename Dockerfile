FROM golang:latest

ADD --keep-git-dir=true  https://github.com/liteseed/edge /repo
WORKDIR /repo
CMD ["make", "docker"]

WORKDIR /
CP /repo/config.example.json /config.json
CP ["/build/docker/edge", /edge]
RUN rm -rf /repo

CMD ["./edge", "generate"]
CMD ["./edge", "start"]
EXPOSE 8080
