FROM golang:alpine AS builder

WORKDIR /buildapp

COPY go.mod go.sum ./

RUN go mod download 

COPY . .

RUN go build -o output ./cmd

FROM alpine:latest

WORKDIR /app 

COPY config.yaml .

COPY --from=builder /buildapp/output .

EXPOSE 8080

CMD ["./output"]

