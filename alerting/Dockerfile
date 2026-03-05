FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o alerting ./alerting/

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/alerting .
COPY --from=builder /app/alerting/rules.yaml .
CMD ["./alerting"]
