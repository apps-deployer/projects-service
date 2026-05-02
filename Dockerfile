FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o projects_service ./cmd/projects_service
RUN CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/migrator

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/projects_service .
COPY --from=builder /build/migrator .
COPY migrations/ migrations/

EXPOSE 50051

CMD ["./projects_service"]
