package grpcclient

import (
	"context"
	"csu-lessons/product/api"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"sync"
)

type Client struct {
	conns *sync.Pool
}

func NewClient(config Config) Client {
	pool := &sync.Pool{
		New: func() interface{} {
			conn, err := grpc.NewClient(
				fmt.Sprintf("%s:%d", config.Host, config.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))

			if err != nil {
				slog.Error(err.Error())
			}

			return conn
		},
	}

	return Client{conns: pool}
}

func (c Client) CheckStock(ctx context.Context, products []int) ([]ProductState, error) {
	client := c.conns.Get().(*grpc.ClientConn)
	defer c.conns.Put(client)

	request := api.CheckStockRequest{}
	for _, product := range products {
		request.ProductIds = append(request.GetProductIds(), int64(product))
	}

	resp, err := api.NewProductServiceClient(client).CheckStock(context.Background(), &request)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return fromProtoModel(resp.ItemsState), nil
}
