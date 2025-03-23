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
	"strings"
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
	ctx context.Context, req *requester.TopLongMemesHighlightRequest,
) (*requester.HighlightResponse, error) {
	memes, err := s.MemeClient.GetTopLongMemes(ctx, req.Limit)
	if err != nil {
		slog.Error("getting top long memes error", "error", err)
		return nil, err
	}

	return &requester.HighlightResponse{Text: formatMemesResponse(memes)}, nil
}

// SearchMemesBySubstring ищет мемы по подстроке
func (s *GrpcServer) SearchMemesBySubstring(
	ctx context.Context, req *requester.SearchHighlightRequest,
) (*requester.HighlightResponse, error) {
	memes, err := s.MemeClient.SearchMemesBySubstring(ctx, req.Query)
	if err != nil {
		slog.Error("getting memes by substring error", "error", err)
		return nil, err
	}

	return &requester.HighlightResponse{Text: formatMemesResponse(memes)}, nil
}

// GetMemesByMonth возвращает мемы за указанный месяц
func (s *GrpcServer) GetMemesByMonth(
	ctx context.Context, req *requester.MonthHighlightRequest,
) (*requester.HighlightResponse, error) {
	memes, err := s.MemeClient.GetMemesByMonth(ctx, req.Month)
	if err != nil {
		slog.Error("getting memes by month error", "error", err)
		return nil, err
	}

	return &requester.HighlightResponse{Text: formatMemesResponse(memes)}, nil
}

// GetRandomMeme возвращает случайный мем
func (s *GrpcServer) GetRandomMeme(
	ctx context.Context, req *requester.EmptyHighlightRequest,
) (*requester.HighlightResponse, error) {
	meme, err := s.MemeClient.GetRandomMeme(ctx)
	if err != nil {
		slog.Error("getting random meme error", "error", err)
		return nil, err
	}

	memes := &search.MemesResponse{Memes: []*search.MemeResponse{meme}}
	return &requester.HighlightResponse{Text: formatMemesResponse(memes)}, nil
}

func formatMemesResponse(memes *search.MemesResponse) string {
	var texts []string
	for _, meme := range memes.Memes {
		texts = append(texts, meme.Text)
	}
	return strings.Join(texts, "\n\n")
}
