FROM golang:latest

RUN mkdir /api

ADD . /api
COPY config.yml /api/cmd/

WORKDIR /api/cmd

RUN go mod download && \
    go build -o hbl .

CMD ["/api/cmd/hbl"]
