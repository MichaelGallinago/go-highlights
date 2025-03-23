FROM golang:latest AS builder

WORKDIR /app

copy go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o requesterService ./cmd/main.go

FROM alpine:latest

COPY --from=builder /app/requesterService /requesterService
COPY --from=builder /app/config.yml /config.yml

EXPOSE 50053

CMD ["./requesterService"]
