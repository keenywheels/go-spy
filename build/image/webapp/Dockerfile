FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk update --no-cache && apk add --no-cache tzdata

COPY go.mod go.sum /webapp/

WORKDIR /webapp

RUN go mod download
COPY . /webapp
RUN go build -ldflags="-s -w" -o /webapp/bin/webapp ./cmd/webapp/main.go

FROM alpine:latest

ENV TZ=Europe/Moscow

EXPOSE 8008

WORKDIR /webapp

COPY --from=builder /usr/share/zoneinfo/${TZ} /etc/localtime
RUN echo $TZ > /etc/timezone

COPY --from=builder /webapp/bin/webapp .
