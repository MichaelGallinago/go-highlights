package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"repositoryService/internal/core/useCase"
	"repositoryService/internal/lib/postgresclient"
	"repositoryService/internal/lib/publishgrpcserver"
	"repositoryService/internal/lib/searchgrpcserver"
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
	PublishGrpcServer publishgrpcserver.Config `yaml:"repository_service_publish"`
	SearchGrpcServer  searchgrpcserver.Config  `yaml:"repository_service_search"`
	PostgresClient    postgresclient.Config    `yaml:"postgresql"`
}

func main() {
	cfg, err := loadConfig("config.yml")
	if err != nil {
		log.Fatalf("Configuration loading error: %v", err)
	}

	postgresclient.InitPostgresClient(cfg.PostgresClient)

	uc := useCase.NewUseCase(
		searchgrpcserver.NewSearchGrpcServer(cfg.SearchGrpcServer),
		publishgrpcserver.NewPublishGrpcServer(cfg.PublishGrpcServer),
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	defer uc.CloseDB()

	slog.Info("Service started...")

	<-sigChan
	slog.Info("The completion signal is received, turning off the service...")

	slog.Info("The service has been stopped")
}
