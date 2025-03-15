package usecase

import (
	"context"
	"csu-lessons/internal/core/entity"
	"csu-lessons/internal/lib/grpcclient"
)

type ProductClient interface {
	CheckStock(ctx context.Context, items []int) ([]grpcclient.ProductState, error)
}

type NotificationClient interface {
	Publish(ctx context.Context, order entity.OrderEvent) error
}
