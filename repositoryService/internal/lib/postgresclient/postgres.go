package postgresclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"log"
)

type PostgresClient struct {
	Pool *pgxpool.Pool
}

func NewPostgresClient(config Config) PostgresClient {
	conString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DbName, config.SslMode)

	pool, err := pgxpool.New(context.Background(), conString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = runMigrations(pool)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	return PostgresClient{Pool: pool}
}

func runMigrations(pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	db := stdlib.OpenDBFromPool(pool)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file:///db/migrations", "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (db *PostgresClient) Close() {
	db.Pool.Close()
}
