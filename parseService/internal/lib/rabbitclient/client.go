package rabbitclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"parseService/internal/core/entity"
	"sync"
)

type RabbitClient struct {
	pool    *sync.Pool
	channel string
}

func NewRabbitClient(config Config) RabbitClient {
	pool := &sync.Pool{
		New: func() interface{} {
			conn, err := amqp091.Dial(
				fmt.Sprintf("amqp://%s:%s@%s:%s",
					config.Login,
					config.Password,
					config.Host,
					config.Port),
			)
			if err != nil {
				slog.Error("Ошибка подключения к RabbitMQ: " + err.Error())
				return nil
			}
			return conn
		},
	}

	return RabbitClient{pool: pool, channel: config.Channel}
}

func (r RabbitClient) GetMessage(ctx context.Context) (entity.Message, error) {
	conn := r.pool.Get().(*amqp091.Connection)
	if conn == nil {
		return entity.Message{}, fmt.Errorf("нет доступных соединений с RabbitMQ")
	}
	defer r.pool.Put(conn)

	ch, err := conn.Channel()
	if err != nil {
		return entity.Message{}, err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		r.channel,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return entity.Message{}, err
	}

	select {
	case msg := <-msgs:
		var message entity.MessageEvent
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			return entity.Message{}, fmt.Errorf("ошибка парсинга JSON: %w", err)
		}
		return entity.Message{
			Timestamp: message.Timestamp,
			Text:      message.Text,
		}, nil
	case <-ctx.Done():
		return entity.Message{}, ctx.Err()
	}
}
