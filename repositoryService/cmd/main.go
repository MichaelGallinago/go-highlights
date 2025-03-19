package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"net"
	"repositoryService/repository"
	"time"

	"google.golang.org/grpc"
)

const connection = "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"

type RepositoryServiceServer struct {
	repository.UnimplementedRepositoryServiceServer
	db *pgxpool.Pool
}

func startGRPCServer(db *pgxpool.Pool) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	repository.RegisterRepositoryServiceServer(grpcServer, &RepositoryServiceServer{db: db})

	slog.Info("gRPC server started on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func connectToDatabase() *pgxpool.Pool {
	conString := fmt.Sprintf(connection, "postgres", "5432", "postgres", "postgres", "postgres", "disable")
	conn, err := pgxpool.New(context.Background(), conString)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	createMemesTable(conn)
	return conn
}

func createMemesTable(conn *pgxpool.Pool) {
	query := `
	CREATE TABLE IF NOT EXISTS memes (
	    id SERIAL PRIMARY KEY,
	    timestamp BIGINT NOT NULL,
	    text TEXT NOT NULL
	);`
	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
}

func main() {
	db := connectToDatabase()
	startGRPCServer(db)
}

func (s *RepositoryServiceServer) PublishMeme(
	ctx context.Context,
	req *repository.PublishMemeRequest,
) (*repository.PublishMemeResponse, error) {
	slog.Info("new meme:", "timestamp", req.Timestamp, "text", req.Text)

	timestamp, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		slog.Error("timestamp parsing error", "error", err)
		return &repository.PublishMemeResponse{Success: false}, err
	}

	query := "INSERT INTO memes (timestamp, text) VALUES ($1, $2)"
	_, err = s.db.Exec(ctx, query, timestamp.Unix(), req.Text)
	if err != nil {
		slog.Error("insert error", "error", err)
		return &repository.PublishMemeResponse{Success: false}, err
	}

	return &repository.PublishMemeResponse{Success: true}, nil
}

func (s *RepositoryServiceServer) GetTopLongMemes(
	ctx context.Context,
	req *repository.TopLongMemesRequest,
) (*repository.MemesResponse, error) {
	query := "SELECT timestamp, text FROM memes ORDER BY LENGTH(text) DESC LIMIT $1"
	rows, err := s.db.Query(ctx, query, req.Limit)
	if err != nil {
		slog.Error("getting top long memes error", "error", err)
		return nil, err
	}
	defer rows.Close()

	var memes []*repository.MemeResponse
	for rows.Next() {
		var timestamp int64
		var text string
		if err := rows.Scan(&timestamp, &text); err != nil {
			return nil, err
		}
		memes = append(memes, &repository.MemeResponse{
			Text:      text,
			Timestamp: fmt.Sprintf("%d", timestamp),
		})
	}
	return &repository.MemesResponse{Memes: memes}, nil
}

func (s *RepositoryServiceServer) SearchMemesBySubstring(
	ctx context.Context,
	req *repository.SearchRequest,
) (*repository.MemesResponse, error) {
	query := "SELECT timestamp, text FROM memes WHERE text ILIKE '%' || $1 || '%'"
	rows, err := s.db.Query(ctx, query, req.Query)
	if err != nil {
		slog.Error("finding meme error", "error", err)
		return nil, err
	}
	defer rows.Close()

	var memes []*repository.MemeResponse
	for rows.Next() {
		var timestamp int64
		var text string
		if err := rows.Scan(&timestamp, &text); err != nil {
			return nil, err
		}
		memes = append(memes, &repository.MemeResponse{
			Text:      text,
			Timestamp: fmt.Sprintf("%d", timestamp),
		})
	}
	return &repository.MemesResponse{Memes: memes}, nil
}

func (s *RepositoryServiceServer) GetMemesByMonth(
	ctx context.Context,
	req *repository.MonthRequest,
) (*repository.MemesResponse, error) {
	year := int(req.Year)
	month := time.Month(req.Month)
	startTime := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).Unix()
	endTime := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC).Unix()

	query := "SELECT timestamp, text FROM memes WHERE timestamp >= $1 AND timestamp < $2"
	rows, err := s.db.Query(ctx, query, startTime, endTime)
	if err != nil {
		slog.Error("getting memes by moths error", "error", err)
		return nil, err
	}
	defer rows.Close()

	var memes []*repository.MemeResponse
	for rows.Next() {
		var timestamp int64
		var text string
		if err := rows.Scan(&timestamp, &text); err != nil {
			return nil, err
		}
		memes = append(memes, &repository.MemeResponse{
			Text:      text,
			Timestamp: fmt.Sprintf("%d", timestamp),
		})
	}
	return &repository.MemesResponse{Memes: memes}, nil
}

func (s *RepositoryServiceServer) GetRandomMeme(
	ctx context.Context,
	req *repository.Empty,
) (*repository.MemeResponse, error) {
	query := "SELECT timestamp, text FROM memes ORDER BY RANDOM() LIMIT 1"
	row := s.db.QueryRow(ctx, query)

	var timestamp int64
	var text string
	err := row.Scan(&timestamp, &text)
	if err != nil {
		slog.Error("getting random meme error", "error", err)
		return nil, err
	}

	return &repository.MemeResponse{
		Text:      text,
		Timestamp: fmt.Sprintf("%d", timestamp),
	}, nil
}
