FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
WORKDIR /app
RUN mkdir data
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
EXPOSE 8080
CMD ["./main"]
