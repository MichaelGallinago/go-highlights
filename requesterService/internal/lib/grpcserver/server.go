package grpcserver

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
	"requesterService/internal/lib/grpcclient"
	"requesterService/repository/api"
	"requesterService/requester"
	"strconv"
)

type GrpcServer struct {
	requester.UnimplementedRequesterServiceServer
	MemeClient *grpcclient.GrpcMemeClient
}

func NewSearchGrpcServer(config Config, memeClient *grpcclient.GrpcMemeClient) *GrpcServer {
	port := strconv.Itoa(config.Port)
	listener, err := net.Listen(config.Network, ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := &GrpcServer{MemeClient: memeClient}
	requester.RegisterRequesterServiceServer(grpcServer, server)

	slog.Info("Requester gRPC server started on port " + port)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return server
}

// GetTopLongMemes возвращает мемы с самой большой длиной
func (s *GrpcServer) GetTopLongMemes(
	ctx context.Context, req *requester.TopLongMemesRequest,
) (*requester.MemesResponse, error) {
	memes, err := s.MemeClient.GetTopLongMemes(ctx, req.Limit)
	if err != nil {
		slog.Error("getting top long memes error", "error", err)
		return nil, err
	}

	return &requester.MemesResponse{Memes: convertMemeResponses(memes.Memes)}, nil
}

// SearchMemesBySubstring ищет мемы по подстроке
func (s *GrpcServer) SearchMemesBySubstring(
	ctx context.Context, req *requester.SearchRequest,
) (*requester.MemesResponse, error) {
	memes, err := s.MemeClient.SearchMemesBySubstring(ctx, req.Query)
	if err != nil {
		slog.Error("getting memes by substring error", "error", err)
		return nil, err
	}

	return &requester.MemesResponse{Memes: convertMemeResponses(memes.Memes)}, nil
}

// GetMemesByMonth возвращает мемы за указанный месяц
func (s *GrpcServer) GetMemesByMonth(
	ctx context.Context, req *requester.MonthRequest,
) (*requester.MemesResponse, error) {
	memes, err := s.MemeClient.GetMemesByMonth(ctx, req.Year, req.Month)
	if err != nil {
		slog.Error("getting memes by month error", "error", err)
		return nil, err
	}

	return &requester.MemesResponse{Memes: convertMemeResponses(memes.Memes)}, nil
}

// GetRandomMeme возвращает случайный мем
func (s *GrpcServer) GetRandomMeme(
	ctx context.Context, req *requester.Empty,
) (*requester.MemeResponse, error) {
	meme, err := s.MemeClient.GetRandomMeme(ctx)
	if err != nil {
		slog.Error("getting random meme error", "error", err)
		return nil, err
	}

	return &requester.MemeResponse{
		Text:      meme.Text,
		Timestamp: meme.Timestamp,
	}, nil
}

func convertMemeResponses(source []*search.MemeResponse) []*requester.MemeResponse {
	target := make([]*requester.MemeResponse, len(source))
	for i, meme := range source {
		if meme == nil {
			continue
		}

		target[i] = &requester.MemeResponse{
			Timestamp: meme.Timestamp,
			Text:      meme.Text,
		}
	}
	return target
}
