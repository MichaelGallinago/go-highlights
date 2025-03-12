package grpc

import (
	"context"
	"csu-lessons/internal/lib/grpcclient"
	"fmt"
)

type productRepository struct {
	client ProductClient
}

func NewProductRepository(client ProductClient) *productRepository {
	return &productRepository{client: client}
}

func (r *productRepository) CheckStock(ctx context.Context, items []int) ([]grpcclient.ProductState, error) {
	states, err := r.client.CheckStock(ctx, items)

	if err != nil {
		if notFounded := findNotFoundItems(items, states); notFounded != nil {
			return states, fmt.Errorf("product stock check warning: %w: %v", ErrNotFound, notFounded)
		}

		if len(states) == 0 {
			return nil, fmt.Errorf("product stock check warning: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return states, nil
}

func findNotFoundItems(items []int, states []grpcclient.ProductState) []int {
	notFoundItems := make([]int, 0)
	for _, itemId := range items {
		founded := false
		for _, state := range states {
			if state.ProductId == int64(itemId) {
				founded = true
			}
		}

		if !founded {
			notFoundItems = append(notFoundItems, itemId)
		}
	}

	return notFoundItems
}
