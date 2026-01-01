FROM golang:1.25-alpine AS builder

WORKDIR /usr/src/starter-api

RUN apk add --no-cache gcc musl-dev

ENV CGO_ENABLED=1

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/app ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /usr/src/starter-api/bin/app ./app

# Copy assets if needed
COPY --from=builder /usr/src/starter-api/assets ./assets
COPY --from=builder /usr/src/starter-api/templates ./templates

# Ensure non-root user owns application files
RUN chown -R appuser:appgroup /app

USER appuser
EXPOSE 3000

CMD ["./app"]