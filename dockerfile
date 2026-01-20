FROM golang:1.25-alpine AS deps
RUN apk --no-cache add bash git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM deps AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/server

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /main
COPY --from=builder /app/migrations /migrations

USER 1001:1001
EXPOSE 8080

CMD ["./main"]