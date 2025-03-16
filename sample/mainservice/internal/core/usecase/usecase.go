package usecase

import (
	"context"
	"csu-lessons/internal/core/entity"
	"log/slog"
	"math/rand"
	"time"
)

type Usecase struct {
	productClient      ProductClient
	notificationClient NotificationClient
}

func NewUsecase(productClient ProductClient, notificationClient NotificationClient) Usecase {
	return Usecase{
		productClient:      productClient,
		notificationClient: notificationClient,
	}
}

func (uc Usecase) CreateOrder(ctx context.Context, order entity.Order) error {
	items := order.GetOrderItems()

	states, err := uc.productClient.CheckStock(ctx, items)
	if err != nil {
		slog.Error("get products state error", err)

		return err
	}

	for _, state := range states {
		if state.Count == 0 {
			slog.Warn("product stock is zero", state.ProductId)
		}
	}

	// create order in other system

	rand.Seed(time.Now().UnixNano())

	randomInt := rand.Intn(1000) + 1
	event := entity.OrderEvent{
		Id:     int64(randomInt),
		UserId: order.UserId,
		Status: "created",
	}
	err = uc.notificationClient.Publish(ctx, event)

	if err != nil {
		slog.Error("send notification error", err)

		return err
	}

	return nil
}
