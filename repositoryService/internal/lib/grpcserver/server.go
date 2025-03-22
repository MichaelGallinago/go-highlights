package grpcserver

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"repositoryService/internal/lib/postgresclient"
	"repositoryService/repository"
	"strconv"
)

type GrpcServer struct {
	repository.UnimplementedRepositoryServiceServer
	DB *postgresclient.PostgresClient
}

// NewGRPCServer запускает gRPC сервер
func NewGRPCServer(config Config, db *postgresclient.PostgresClient) *GrpcServer {
	port := strconv.Itoa(config.Port)
	listener, err := net.Listen(config.Network, ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := &GrpcServer{DB: db}
	repository.RegisterRepositoryServiceServer(grpcServer, server)

	log.Println("gRPC server started on port " + port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return server
}
