package httpengine

import (
	"context"
	"csu-lessons/internal/delivery/httpengine/handler"
	"github.com/gin-gonic/gin"
)

func InitRouter(ctx context.Context, config Config, handler *handler.Handler) *gin.Engine {
	router := gin.Default()
	router.POST("/create-order", handler.CreateOrder)
	return router
}
