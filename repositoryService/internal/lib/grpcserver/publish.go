package grpcserver

import (
	"context"
	"log/slog"
	"repositoryService/repository"
	"time"
)

// PublishMeme сохраняет мем в БД
func (s *GrpcServer) PublishMeme(
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
