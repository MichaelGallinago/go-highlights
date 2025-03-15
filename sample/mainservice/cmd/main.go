package main

import (
	"context"
	"csu-lessons/internal/core/usecase"
	"csu-lessons/internal/delivery/httpengine"
	"csu-lessons/internal/delivery/httpengine/handler"
	"csu-lessons/internal/lib/grpcclient"
	"csu-lessons/internal/lib/rabbitclient"
	"csu-lessons/internal/repository/datasources/grpc"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log/slog"
	"net/http"
	"os"
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
	ProductService grpcclient.Config   `yaml:"product_service"`
	Router         httpengine.Config   `yaml:"router"`
	Rabbit         rabbitclient.Config `yaml:"rabbit"`
}

func main() {
	gin.SetMode(gin.DebugMode)

	cfg, err := loadConfig("config.yml")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	productClient := grpcclient.NewClient(cfg.ProductService)

	repo := grpc.NewProductRepository(productClient)

	notificationClient := rabbitclient.NewRabbitClient(cfg.Rabbit)

	uc := usecase.NewUsecase(repo, notificationClient)

	productHandler := handler.NewHandler(uc)

	router := httpengine.InitRouter(context.Background(), cfg.Router, productHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Router.Port), router); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
