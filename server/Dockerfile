FROM golang:1.24-bookworm AS build

WORKDIR /build
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=1 go build -tags=viper_bind_struct -o /build/til ./cmd

FROM debian:bookworm AS run

RUN apt update ; apt install -y ca-certificates
WORKDIR /
COPY --from=build /build/til /til
RUN chmod +x /til
ENTRYPOINT ["/til"]