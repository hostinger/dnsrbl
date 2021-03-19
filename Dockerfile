# Compile hbl binary
FROM golang:1.16.2 AS build-env

RUN mkdir /api

COPY . /api

WORKDIR /api/cmd/hbl

RUN go mod download && \
    go build -o hbl .

# Create final image
FROM gcr.io/distroless/base:nonroot

COPY --from=build-env /api/cmd/hbl/hbl /api/cmd/hbl/hbl

USER nobody

CMD ["/api/cmd/hbl/hbl"]
