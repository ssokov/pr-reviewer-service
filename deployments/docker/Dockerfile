FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o pr-reviewer-service ./cmd/pr-reviewer-service

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/pr-reviewer-service .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./pr-reviewer-service"]