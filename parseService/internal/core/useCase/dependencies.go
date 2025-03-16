package useCase

import (
	"context"
	"parseService/internal/core/entity"
)

type MemeClient interface {
	Publish(ctx context.Context, meme entity.Meme) error
}

type NotificationClient interface {
	GetMessage(ctx context.Context) (entity.Message, error)
}
