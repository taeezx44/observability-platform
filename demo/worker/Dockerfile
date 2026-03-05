FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o demo-worker

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/

COPY --from=builder /app/demo-worker .
COPY --from=builder /app/go.mod go.sum ./

EXPOSE 8082

CMD ["./demo-worker"]
