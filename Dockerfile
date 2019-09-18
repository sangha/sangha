FROM golang:alpine AS builder

LABEL authors="Christian Muehlhaeuser: muesli@gmail.com"

# Install git & make
# Git is required for fetching the dependencies
RUN apk update && \
    apk add --no-cache git make ca-certificates && \
    update-ca-certificates

# Set the working directory for the container
WORKDIR /go/sangha

# Build the binary
COPY . .
RUN go build

FROM alpine

RUN apk update && \
    apk add --no-cache ca-certificates tzdata && \
    update-ca-certificates

COPY --from=builder /go/sangha/sangha /go/bin/sangha

# Expose the application port
EXPOSE 9980

# create a volume for the configuration persistence
VOLUME /conf

# This form of ENTRYPOINT allows the sangha process to catch signals from the `docker stop` command
ENTRYPOINT /go/bin/sangha --config /conf/sangha.conf serve
