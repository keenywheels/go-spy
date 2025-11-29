FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk update --no-cache && apk add --no-cache tzdata

COPY go.mod go.sum /scheduler/

WORKDIR /scheduler

RUN go mod download
COPY . /scheduler
RUN go build -ldflags="-s -w" -o /scheduler/bin/scheduler ./cmd/scheduler/main.go

FROM alpine:latest

ENV TZ=Europe/Moscow

WORKDIR /scheduler

COPY --from=builder /usr/share/zoneinfo/${TZ} /etc/localtime
RUN echo $TZ > /etc/timezone

COPY --from=builder /scheduler/bin/scheduler .
