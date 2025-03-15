package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"log"
	"net"
	api "productservice/product"
)

const connection = "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	conString := fmt.Sprintf(connection, "postgres", "5432", "postgres", "postgres", "postgres", "disable")
	conn, err := pgxpool.New(context.Background(), conString)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	s := grpc.NewServer()
	api.RegisterProductServiceServer(s, &Server{db: conn})
	log.Println("Starting server...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type Server struct {
	api.UnimplementedProductServiceServer
	db *pgxpool.Pool
}

func (s Server) CheckStock(ctx context.Context, request *api.CheckStockRequest) (*api.CheckStockResponse, error) {
	stock := map[int64]int64{}

	rows, err := s.db.Query(ctx, `SELECT product_id, count FROM products WHERE product_id = ANY ($1::BIGINT[])`, request.ProductIds)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var productId int64
		var count int64
		if err := rows.Scan(&productId, &count); err != nil {
			return nil, err
		}
		stock[productId] = count
	}

	var response api.CheckStockResponse

	for _, id := range request.ProductIds {
		response.ItemsState = append(response.ItemsState, &api.ProductState{
			ProductId: id,
			Count:     stock[id],
		})
	}

	return &response, nil
}

func (s Server) mustEmbedUnimplementedProductServiceServer() {
}
