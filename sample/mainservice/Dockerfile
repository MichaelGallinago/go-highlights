FROM golang:latest AS builder

WORKDIR /app

copy go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o http-gateway ./cmd

FROM alpine:latest

COPY --from=builder /app/config.yml /config.yml
COPY --from=builder /app/http-gateway /http-gateway

EXPOSE 8021

CMD ["/http-gateway"]