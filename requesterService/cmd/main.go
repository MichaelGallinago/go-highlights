package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"os/signal"
	"parseService/internal/core/useCase"
	"parseService/internal/lib/grpcclient"
	"parseService/internal/lib/rabbitclient"
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
	RepositoryService grpcclient.Config   `yaml:"repository_service"`
	Rabbit            rabbitclient.Config `yaml:"rabbit"`
}

func main() {
	gin.SetMode(gin.DebugMode)

	cfg, err := loadConfig("config.yml")
	if err != nil {
		slog.Error("Ошибка загрузки конфигурации:", err)
		os.Exit(1)
	}

	memeClient := grpcclient.NewGrpcMemeClient(cfg.RepositoryService)
	notificationServer := rabbitclient.NewRabbitClient(cfg.Rabbit)
	uc := useCase.NewUseCase(memeClient, notificationClient)

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := uc.StartMessageListening(ctx); err != nil {
			slog.Error("Ошибка во время прослушивания сообщений:", err)
		}
	}()

	slog.Info("Сервис запущен и слушает сообщения...")

	<-sigChan
	slog.Info("Получен сигнал завершения, выключаем сервис...")

	cancel()

	slog.Info("Сервис остановлен")
}
