package handler

import (
	"context"
	"csu-lessons/internal/core/entity"
)

type OrderClient interface {
	CreateOrder(ctx context.Context, order entity.Order) error
}
