package useCase

import (
	"context"
	"fmt"
	"log"
	"parseService/internal/core/entity"
	"regexp"
	"strings"
	"time"
)

type UseCase struct {
	memeClient         MemeClient
	notificationClient NotificationClient
}

func NewUseCase(memeClient MemeClient, notificationClient NotificationClient) UseCase {
	return UseCase{
		memeClient:         memeClient,
		notificationClient: notificationClient,
	}
}

func (uc UseCase) StartMessageListening(ctx context.Context) error {
	log.Println("Начало прослушивания сообщений...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Остановка прослушивания сообщений.")
			return ctx.Err()
		default:

			message, err := uc.notificationClient.GetMessage(ctx)
			if err != nil {
				log.Println("Ошибка при получении сообщения:", err)
				time.Sleep(2 * time.Second)
				continue
			}

			log.Println("Получено сообщение:", message.Text)

			memes := parseMemes(message.Text)

			for _, meme := range memes {
				err := uc.memeClient.Publish(ctx, entity.Meme{
					Timestamp: message.Timestamp,
					Text:      meme,
				})
				if err != nil {
					log.Println("Ошибка при публикации мема:", err)
				} else {
					log.Println("Мем опубликован:", meme)
				}
			}
		}
	}
}

func parseMemes(text string) []string {
	var memes []string

	re := regexp.MustCompile(`(?m)"((?:[^"\\]|\\.)+)"\s*(?:\((с|С)\)|©)?\s*([A-Za-zА-Яа-яёЁ, ]*)`)
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		quote := match[1]
		author := match[3]

		var formatted string

		if author == "" {
			formatted = fmt.Sprintf("\"%s\"", quote)
		} else {
			formatted = fmt.Sprintf("\"%s\" (%s)", quote, author)
		}

		formatted = strings.ReplaceAll(formatted, `\"`, `"`)
		memes = append(memes, formatted)
	}

	return memes
}
