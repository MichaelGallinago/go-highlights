package publishgrpcserver

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
	"repositoryService/internal/lib/postgresclient"
	"repositoryService/repository"
	"strconv"
	"time"
)

type PublishGrpcServer struct {
	repository.UnimplementedRepositoryServiceServer
	DB postgresclient.PostgresClient
}

// NewPublishGrpcServer запускает gRPC сервер
func NewPublishGrpcServer(config Config, db postgresclient.PostgresClient) *PublishGrpcServer {
	port := strconv.Itoa(config.Port)
	listener, err := net.Listen(config.Network, ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := &PublishGrpcServer{DB: db}

	repository.RegisterRepositoryServiceServer(grpcServer, server)

	log.Println("Publish gRPC server started on port " + port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return server
}

// PublishMeme сохраняет мем в БД
func (s *PublishGrpcServer) PublishMeme(
	ctx context.Context, req *repository.PublishMemeRequest,
) (*repository.PublishMemeResponse, error) {
	timestamp, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		slog.Error("timestamp parsing error", "error", err)
		return &repository.PublishMemeResponse{Success: false}, err
	}

	err = s.DB.InsertMeme(ctx, timestamp.Unix(), req.Text)
	if err != nil {
		slog.Error("insert error", "error", err)
		return &repository.PublishMemeResponse{Success: false}, err
	}

	slog.Info("new meme inserted:", "timestamp", req.Timestamp, "text", req.Text)
	return &repository.PublishMemeResponse{Success: true}, nil
}
