FROM golang:1.22-alpine AS builder

WORKDIR /usr/src/starter-api

RUN apk add --no-cache gcc musl-dev

ENV CGO_ENABLED=1

COPY go.mod go.sum ./
RUN go mod download
COPY . .

EXPOSE 3000