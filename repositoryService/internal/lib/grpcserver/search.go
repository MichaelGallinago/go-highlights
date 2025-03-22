package grpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"repositoryService/internal/core/entity"
	"repositoryService/repository"
	"time"
)

// GetTopLongMemes возвращает мемы с самой большой длиной
func (s *GrpcServer) GetTopLongMemes(
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
func (s *GrpcServer) SearchMemesBySubstring(
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
func (s *GrpcServer) GetMemesByMonth(
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
func (s *GrpcServer) GetRandomMeme(
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
