package useCase

import (
	"context"
	"fmt"
	"log"
	"parseService/internal/core/entity"
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

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Остановка прослушивания сообщений.")
			return ctx.Err()

		case <-ticker.C:
			message, err := uc.notificationClient.GetMessage(ctx)
			if err != nil {
				log.Println("Ошибка при получении сообщения:", err)
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
	splitPoint := strings.Index(text, "\"")
	if splitPoint != -1 {
		text = text[splitPoint:]
	}

	blocks := strings.Split(text, "\n\n")

	var memes []string

	for _, block := range blocks {
		quote, author := parseQuoteAndAuthor(block)

		if quote == "" {
			continue
		}

		if author != "" {
			memes = append(memes, fmt.Sprintf("%s (%s)", quote, author))
		} else {
			memes = append(memes, fmt.Sprintf("%s", quote))
		}
	}

	return memes
}

func parseQuoteAndAuthor(quote string) (resultQuote, author string) {
	resultQuote = strings.TrimSpace(quote)

	lastLineStartIndex := strings.LastIndex(resultQuote, "\n")
	author = ""

	if lastLineStartIndex == -1 {
		return resultQuote, author
	}

	lastLine := strings.TrimSpace(resultQuote[lastLineStartIndex:])

	switch {
	case strings.HasPrefix(lastLine, "(с)") || strings.HasPrefix(lastLine, "(С)"):
		author = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(lastLine, "(С)"), "(с)"))
		resultQuote = resultQuote[:lastLineStartIndex]
	case strings.HasPrefix(lastLine, "©"):
		author = strings.TrimSpace(strings.TrimPrefix(lastLine, "©"))
		resultQuote = resultQuote[:lastLineStartIndex]
	case !strings.ContainsAny(lastLine, `"„«»`):
		author = lastLine
		resultQuote = resultQuote[:lastLineStartIndex]
	}

	return resultQuote, author
}
