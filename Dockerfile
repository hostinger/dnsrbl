# Compile hbl binary
FROM golang:1.16.2 AS build-env

RUN mkdir /api

COPY . /api

WORKDIR /api

#Generate docs
RUN go get -u github.com/swaggo/swag/cmd/swag && \
    go get -u github.com/alecthomas/template && \
    swag init -g cmd/hbl/hbl.go

WORKDIR /api/cmd/hbl

#Build api
RUN go mod download && \
    go build -o hbl .

# Create final image
FROM gcr.io/distroless/base:nonroot

COPY --from=build-env /api/cmd/hbl/hbl /api/cmd/hbl/hbl

CMD ["/api/cmd/hbl/hbl"]
