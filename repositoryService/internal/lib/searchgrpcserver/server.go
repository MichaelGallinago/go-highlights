package searchgrpcserver

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
	"repositoryService/internal/core/entity"
	"repositoryService/internal/lib/postgresclient"
	"repositoryService/repository/search"
	"strconv"
)

type SearchGrpcServer struct {
	search.UnimplementedRepositoryServiceSearchServer
	DB postgresclient.PostgresClient
}

func NewSearchGrpcServer(config Config, db postgresclient.PostgresClient) *SearchGrpcServer {
	port := strconv.Itoa(config.Port)
	listener, err := net.Listen(config.Network, ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := &SearchGrpcServer{DB: db}
	search.RegisterRepositoryServiceSearchServer(grpcServer, server)

	slog.Info("Search gRPC server started on port " + port)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return server
}

// GetTopLongMemes возвращает мемы с самой большой длиной
func (s *SearchGrpcServer) GetTopLongMemes(
	ctx context.Context, req *search.TopLongMemesRequest,
) (*search.MemesResponse, error) {
	memes, err := s.DB.GetTopLongMemes(ctx, int(req.Limit))
	if err != nil {
		slog.Error("getting top long memes error", "error", err)
		return nil, err
	}

	return mapMemesToResponse(memes), nil
}

// SearchMemesBySubstring ищет мемы по подстроке
func (s *SearchGrpcServer) SearchMemesBySubstring(
	ctx context.Context, req *search.SearchRequest,
) (*search.MemesResponse, error) {
	memes, err := s.DB.SearchMemes(ctx, req.Query)
	if err != nil {
		slog.Error("finding meme error", "error", err)
		return nil, err
	}

	return mapMemesToResponse(memes), nil
}

func (s *SearchGrpcServer) GetMemesByMonth(
	ctx context.Context, req *search.MonthRequest,
) (*search.MemesResponse, error) {
	memes, err := s.DB.GetMemesByMonth(ctx, req.Month)
	if err != nil {
		slog.Error("getting memes by month error", "error", err)
		return nil, err
	}

	return mapMemesToResponse(memes), nil
}

// GetRandomMeme возвращает случайный мем
func (s *SearchGrpcServer) GetRandomMeme(
	ctx context.Context, req *search.Empty,
) (*search.MemeResponse, error) {
	meme, err := s.DB.GetRandomMeme(ctx)
	if err != nil {
		slog.Error("getting random meme error", "error", err)
		return nil, err
	}

	return &search.MemeResponse{
		Text:      meme.Text,
		Timestamp: fmt.Sprintf("%d", meme.Timestamp),
	}, nil
}

// mapMemesToResponse конвертирует массив мемов в gRPC-ответ
func mapMemesToResponse(memes []entity.Meme) *search.MemesResponse {
	var responses []*search.MemeResponse
	for _, meme := range memes {
		responses = append(responses, &search.MemeResponse{
			Text:      meme.Text,
			Timestamp: fmt.Sprintf("%d", meme.Timestamp),
		})
	}
	return &search.MemesResponse{Memes: responses}
}
