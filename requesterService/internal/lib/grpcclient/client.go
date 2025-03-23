package grpcclient

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"requesterService/repository/api"
	"sync"
)

type GrpcMemeClient struct {
	conns *sync.Pool
}

func NewGrpcMemeClient(config Config) *GrpcMemeClient {
	pool := &sync.Pool{
		New: func() interface{} {
			conn, err := grpc.NewClient(
				fmt.Sprintf("%s:%d", config.Host, config.Port),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				slog.Error("gRPC-connection error: ", err)
			}
			return conn
		},
	}

	return &GrpcMemeClient{conns: pool}
}

// GetTopLongMemes возвращает мемы с самой большой длиной
func (c *GrpcMemeClient) GetTopLongMemes(
	ctx context.Context, limit int32,
) (*search.MemesResponse, error) {
	client := c.conns.Get().(*grpc.ClientConn)
	defer c.conns.Put(client)

	response, err := search.NewRepositoryServiceSearchClient(client).GetTopLongMemes(ctx, &search.TopLongMemesRequest{
		Limit: limit,
	})
	if err != nil {
		slog.Error("Error on GetTopLongMemes: ", err)
	}

	return response, err
}

// SearchMemesBySubstring ищет мемы по подстроке
func (c *GrpcMemeClient) SearchMemesBySubstring(
	ctx context.Context, query string,
) (*search.MemesResponse, error) {
	client := c.conns.Get().(*grpc.ClientConn)
	defer c.conns.Put(client)

	response, err := search.NewRepositoryServiceSearchClient(client).SearchMemesBySubstring(ctx, &search.SearchRequest{
		Query: query,
	})
	if err != nil {
		slog.Error("Error on SearchMemesBySubstring: ", err)
	}

	return response, err
}

// GetMemesByMonth возвращает мемы за указанный месяц
func (c *GrpcMemeClient) GetMemesByMonth(
	ctx context.Context, month int32,
) (*search.MemesResponse, error) {
	client := c.conns.Get().(*grpc.ClientConn)
	defer c.conns.Put(client)

	response, err := search.NewRepositoryServiceSearchClient(client).GetMemesByMonth(ctx, &search.MonthRequest{
		Month: month,
	})
	if err != nil {
		slog.Error("Error on SearchMemesBySubstring: ", err)
	}

	return response, err
}

// GetRandomMeme возвращает случайный мем
func (c *GrpcMemeClient) GetRandomMeme(
	ctx context.Context,
) (*search.MemeResponse, error) {
	client := c.conns.Get().(*grpc.ClientConn)
	defer c.conns.Put(client)

	response, err := search.NewRepositoryServiceSearchClient(client).GetRandomMeme(ctx, &search.Empty{})
	if err != nil {
		slog.Error("Error on GetRandomMeme: ", err)
	}

	return response, err
}
