FROM golang:1.24-bookworm AS build

WORKDIR /build
COPY . .
RUN go mod download

ENTRYPOINT ["go", "run", "./"]