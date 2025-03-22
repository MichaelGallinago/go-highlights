package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"repositoryService/internal/core/useCase"
	"repositoryService/internal/lib/grpcserver"
	"repositoryService/internal/lib/postgresclient"
	"syscall"
)

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

type Config struct {
	GrpcServer     grpcserver.Config     `yaml:"repository_service"`
	PostgresClient postgresclient.Config `yaml:"postgresql"`
}

func main() {
	cfg, err := loadConfig("config.yml")
	if err != nil {
		log.Fatalf("Configuration loading error: %v", err)
	}

	db := postgresclient.NewPostgresClient(cfg.PostgresClient)
	grpcServer := grpcserver.NewGRPCServer(cfg.GrpcServer, db)
	uc := useCase.NewUseCase(grpcServer, grpcServer, db)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	defer uc.CloseDB()

	slog.Info("Service started...")

	<-sigChan
	slog.Info("The completion signal is received, turning off the service...")

	slog.Info("The service has been stopped")
}
