package grpc

import (
	"context"
	"csu-lessons/internal/lib/grpcclient"
)

type ProductClient interface {
	CheckStock(ctx context.Context, products []int) ([]grpcclient.ProductState, error)
}
