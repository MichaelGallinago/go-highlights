FROM golang:latest AS builder

WORKDIR /app

copy go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o consumer ./cmd

FROM alpine:latest

COPY --from=builder /app/consumer /consumer

CMD ["/consumer"]