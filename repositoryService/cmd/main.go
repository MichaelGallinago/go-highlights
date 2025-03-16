package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"repositoryService/repository"

	"google.golang.org/grpc"
)

type MemeServiceServer struct {
	repository.UnimplementedMemeServiceServer
}

func (s *MemeServiceServer) PublishMeme(ctx context.Context, req *repository.MemeRequest) (*repository.MemeResponse, error) {
	slog.Info("Получен мем:", "timestamp", req.Timestamp, "text", req.Text)

	return &repository.MemeResponse{Success: true}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	repository.RegisterMemeServiceServer(grpcServer, &MemeServiceServer{})

	slog.Info("gRPC сервер запущен на порту 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
