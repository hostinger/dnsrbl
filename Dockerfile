FROM golang:latest

RUN mkdir /api

ADD . /api

WORKDIR /api/cmd

RUN go mod download && \
    go build -o hbl .

CMD ["/api/cmd/hbl"]
