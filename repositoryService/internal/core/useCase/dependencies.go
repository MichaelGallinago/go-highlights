package useCase

import (
	"context"
	"repositoryService/internal/core/entity"
	"repositoryService/repository"
)

type SearchServer interface {
	// GetTopLongMemes возвращает мемы с самой большой длиной
	GetTopLongMemes(ctx context.Context, req *repository.TopLongMemesRequest) (*repository.MemesResponse, error)

	// SearchMemesBySubstring ищет мемы по подстроке
	SearchMemesBySubstring(ctx context.Context, req *repository.SearchRequest) (*repository.MemesResponse, error)

	// GetMemesByMonth возвращает мемы за указанный месяц
	GetMemesByMonth(ctx context.Context, req *repository.MonthRequest) (*repository.MemesResponse, error)

	// GetRandomMeme возвращает случайный мем
	GetRandomMeme(ctx context.Context, req *repository.Empty) (*repository.MemeResponse, error)
}

type PublishServer interface {
	// PublishMeme сохраняет мем
	PublishMeme(ctx context.Context, req *repository.PublishMemeRequest) (*repository.PublishMemeResponse, error)
}

type DB interface {
	// InsertMeme добавляет мем в БД
	InsertMeme(ctx context.Context, timestamp int64, text string) error

	// GetTopLongMemes возвращает топ мемов по длине текста
	GetTopLongMemes(ctx context.Context, limit int) ([]entity.Meme, error)

	// SearchMemes ищет мемы по подстроке
	SearchMemes(ctx context.Context, query string) ([]entity.Meme, error)

	// GetMemesByMonth возвращает мемы за указанный период (месяц)
	GetMemesByMonth(ctx context.Context, startTime, endTime int64) ([]entity.Meme, error)

	// GetRandomMeme возвращает случайный мем
	GetRandomMeme(ctx context.Context) (entity.Meme, error)

	Close()
}
