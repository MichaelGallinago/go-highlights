package rabbitclient

import (
	"context"
	"csu-lessons/internal/core/entity"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"sync"
)

type RabbitClient struct {
	pool    *sync.Pool
	channel string
}

func NewRabbitClient(config Config) RabbitClient {
	pool := sync.Pool{
		New: func() interface{} {
			conn, err := amqp091.Dial(
				fmt.Sprintf("amqp://%s:%s@%s:%s",
					config.Login,
					config.Password,
					config.Host,
					config.Port),
			)
			if err != nil {
				slog.Error(err.Error())
			}

			return conn
		},
	}

	return RabbitClient{pool: &pool, channel: config.Channel}
}

func (r RabbitClient) Publish(ctx context.Context, order entity.OrderEvent) error {
	conn := r.pool.Get().(*amqp091.Connection)
	defer r.pool.Put(conn)

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	res, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		r.channel,
		false,
		false,
		amqp091.Publishing{ContentType: "application/json", Body: res})

	return err
}
