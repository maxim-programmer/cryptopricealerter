FROM golang:1.25.4 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cpa ./cmd/alerter

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/cpa .
RUN chmod +x cpa
CMD ["./cpa"]