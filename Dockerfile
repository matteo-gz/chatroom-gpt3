FROM golang:1.19 AS builder
COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

COPY --from=builder /src/bin /app
COPY --from=builder /src/html /app/html

WORKDIR /app
EXPOSE 8000

