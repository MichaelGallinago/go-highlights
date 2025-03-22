package grpcclient

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"parseService/internal/core/entity"
	repository "parseService/repository/api"
	"sync"
)

type GrpcMemeClient struct {
	conns *sync.Pool
}

func NewGrpcMemeClient(config Config) GrpcMemeClient {
	pool := &sync.Pool{
		New: func() interface{} {
			conn, err := grpc.NewClient(
				fmt.Sprintf("%s:%d", config.Host, config.Port),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				slog.Error("Ошибка gRPC-подключения: ", err)
			}
			return conn
		},
	}

	return GrpcMemeClient{conns: pool}
}

func (c GrpcMemeClient) Publish(ctx context.Context, meme entity.Meme) error {
	client := c.conns.Get().(*grpc.ClientConn)
	defer c.conns.Put(client)

	_, err := repository.NewRepositoryServicePublishClient(client).PublishMeme(ctx, &repository.PublishMemeRequest{
		Timestamp: meme.Timestamp.Format("2006-01-02T15:04:05Z"),
		Text:      meme.Text,
	})
	if err != nil {
		slog.Error("Ошибка при отправке мема:", err)
	}

	return err
}
