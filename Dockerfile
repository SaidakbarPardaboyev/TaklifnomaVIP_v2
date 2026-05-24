FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -o warehouse-operation-service ./cmd/main.go

FROM alpine:3.23.2
RUN apk add --no-cache tzdata
ENV TZ=Asia/Tashkent
WORKDIR /app
COPY --from=builder /app/warehouse-operation-service ./
CMD ["./warehouse-operation-service"]