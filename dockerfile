# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /migrate ./cmd/migrate

# Runtime stage

FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app


COPY --from=builder /server .
COPY --from=builder /migrate .


COPY migrations ./migrations
COPY config ./config


RUN adduser -D -s /bin/sh appuser
USER appuser

CMD ["./server"]