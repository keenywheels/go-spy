FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk update --no-cache && apk add --no-cache tzdata

COPY go.mod go.sum /gospy/

WORKDIR /gospy

RUN go mod download
COPY . /gospy
RUN go build -ldflags="-s -w" -o /gospy/bin/gospy ./cmd/gospy/main.go

FROM alpine:latest

ENV TZ=Europe/Moscow

EXPOSE 8008

WORKDIR /gospy

COPY --from=builder /usr/share/zoneinfo/${TZ} /etc/localtime
RUN echo $TZ > /etc/timezone

COPY --from=builder /gospy/bin/gospy .
