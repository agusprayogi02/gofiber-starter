FROM golang:1.25-alpine AS builder

WORKDIR /usr/src/starter-api

RUN apk add --no-cache gcc musl-dev
RUN go install github.com/air-verse/air@latest

ENV CGO_ENABLED=1

COPY go.mod ./
RUN go mod download && go mod verify
COPY . .

EXPOSE 3000