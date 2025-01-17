FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN go build -o /go-auth-service ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /go-auth-service .

EXPOSE 3000

CMD ["./go-auth-service"]
