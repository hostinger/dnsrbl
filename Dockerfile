FROM golang:latest

RUN mkdir /api

ADD . /api

WORKDIR /api/cmd/hbl

RUN go mod download && \
    go build -o hbl .

CMD ["/api/cmd/hbl/hbl"]
