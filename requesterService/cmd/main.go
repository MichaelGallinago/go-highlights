package main

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"os/signal"
	"requesterService/internal/core/useCase"
	"requesterService/internal/lib/grpcclient"
	"requesterService/internal/lib/grpcserver"
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
	RepositoryService grpcclient.Config `yaml:"repository_service"`
	RequesterService  grpcserver.Config `yaml:"requester_service"`
}

func main() {
	cfg, err := loadConfig("config.yml")
	if err != nil {
		slog.Error("Configuration loading error:", err)
		os.Exit(1)
	}

	memeClient := grpcclient.NewGrpcMemeClient(cfg.RepositoryService)
	requesterServer := grpcserver.NewSearchGrpcServer(cfg.RequesterService, memeClient)
	_ = useCase.NewUseCase(memeClient, requesterServer)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	slog.Info("Service started...")

	<-sigChan
	slog.Info("The completion signal is received, turning off the service...")

	slog.Info("The service has been stopped")
}
