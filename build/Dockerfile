
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
RUN usermod -u 1001 -g 1001 appuser
RUN chown -R appuser:appuser /app
USER appuser
EXPOSE 9001
CMD ["./main"]
