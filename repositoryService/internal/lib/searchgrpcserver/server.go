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
	"repositoryService/repository"
	"strconv"
	"time"
)

type SearchGrpcServer struct {
	repository.UnimplementedRepositoryServiceServer
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
	repository.RegisterRepositoryServiceServer(grpcServer, server)

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
	ctx context.Context, req *repository.TopLongMemesRequest,
) (*repository.MemesResponse, error) {
	memes, err := s.DB.GetTopLongMemes(ctx, int(req.Limit))
	if err != nil {
		slog.Error("getting top long memes error", "error", err)
		return nil, err
	}

	return mapMemesToResponse(memes), nil
}

// SearchMemesBySubstring ищет мемы по подстроке
func (s *SearchGrpcServer) SearchMemesBySubstring(
	ctx context.Context, req *repository.SearchRequest,
) (*repository.MemesResponse, error) {
	memes, err := s.DB.SearchMemes(ctx, req.Query)
	if err != nil {
		slog.Error("finding meme error", "error", err)
		return nil, err
	}

	return mapMemesToResponse(memes), nil
}

// GetMemesByMonth возвращает мемы за указанный месяц
func (s *SearchGrpcServer) GetMemesByMonth(
	ctx context.Context, req *repository.MonthRequest,
) (*repository.MemesResponse, error) {
	year := int(req.Year)
	month := time.Month(req.Month)

	startTime := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).Unix()
	endTime := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC).Unix()

	memes, err := s.DB.GetMemesByMonth(ctx, startTime, endTime)
	if err != nil {
		slog.Error("getting memes by month error", "error", err)
		return nil, err
	}

	return mapMemesToResponse(memes), nil
}

// GetRandomMeme возвращает случайный мем
func (s *SearchGrpcServer) GetRandomMeme(
	ctx context.Context, req *repository.Empty,
) (*repository.MemeResponse, error) {
	meme, err := s.DB.GetRandomMeme(ctx)
	if err != nil {
		slog.Error("getting random meme error", "error", err)
		return nil, err
	}

	return &repository.MemeResponse{
		Text:      meme.Text,
		Timestamp: fmt.Sprintf("%d", meme.Timestamp),
	}, nil
}

// mapMemesToResponse конвертирует массив мемов в gRPC-ответ
func mapMemesToResponse(memes []entity.Meme) *repository.MemesResponse {
	var responses []*repository.MemeResponse
	for _, meme := range memes {
		responses = append(responses, &repository.MemeResponse{
			Text:      meme.Text,
			Timestamp: fmt.Sprintf("%d", meme.Timestamp),
		})
	}
	return &repository.MemesResponse{Memes: responses}
}
