FROM golang:latest AS builder

WORKDIR /app

copy go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o grpc-server ./cmd

FROM alpine:latest

COPY --from=builder /app/grpc-server /grpc-server

EXPOSE 50051
EXPOSE 50052

CMD ["/grpc-server"]