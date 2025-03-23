package useCase

import (
	"context"
	"repositoryService/internal/core/entity"
	"repositoryService/repository/publish"
	"repositoryService/repository/search"
)

type SearchServer interface {
	// GetTopLongMemes возвращает мемы с самой большой длиной
	GetTopLongMemes(ctx context.Context, req *search.TopLongMemesRequest) (*search.MemesResponse, error)

	// SearchMemesBySubstring ищет мемы по подстроке
	SearchMemesBySubstring(ctx context.Context, req *search.SearchRequest) (*search.MemesResponse, error)

	// GetMemesByMonth возвращает мемы за указанный месяц
	GetMemesByMonth(ctx context.Context, req *search.MonthRequest) (*search.MemesResponse, error)

	// GetRandomMeme возвращает случайный мем
	GetRandomMeme(ctx context.Context, req *search.Empty) (*search.MemeResponse, error)
}

type PublishServer interface {
	// PublishMeme сохраняет мем
	PublishMeme(ctx context.Context, req *publish.PublishMemeRequest) (*publish.PublishMemeResponse, error)
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
