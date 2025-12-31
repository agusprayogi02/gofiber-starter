FROM golang:1.25-alpine 

WORKDIR /usr/src/starter-api

RUN apk add --no-cache gcc musl-dev

ENV CGO_ENABLED=1

COPY go.mod ./
RUN go mod download && go mod verify
COPY . .

RUN go build -o /usr/bin/starter-api

EXPOSE 3000