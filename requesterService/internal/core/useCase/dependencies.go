package useCase

import (
	"context"
	search "requesterService/repository/api"
	"requesterService/requester"
)

type MemeClient interface {
	// GetTopLongMemes возвращает мемы с самой большой длиной
	GetTopLongMemes(ctx context.Context, limit int32) (*search.MemesResponse, error)

	// SearchMemesBySubstring ищет мемы по подстроке
	SearchMemesBySubstring(ctx context.Context, query string) (*search.MemesResponse, error)

	// GetMemesByMonth возвращает мемы за указанный месяц
	GetMemesByMonth(ctx context.Context, year int32, month int32) (*search.MemesResponse, error)

	// GetRandomMeme возвращает случайный мем
	GetRandomMeme(ctx context.Context) (*search.MemeResponse, error)
}

type RequesterServer interface {
	// GetTopLongMemes возвращает мемы с самой большой длиной
	GetTopLongMemes(ctx context.Context, req *requester.TopLongMemesRequest) (*requester.HighlightResponse, error)

	// SearchMemesBySubstring ищет мемы по подстроке
	SearchMemesBySubstring(ctx context.Context, req *requester.SearchRequest) (*requester.HighlightResponse, error)

	// GetMemesByMonth возвращает мемы за указанный месяц
	GetMemesByMonth(ctx context.Context, req *requester.MonthRequest) (*requester.HighlightResponse, error)

	// GetRandomMeme возвращает случайный мем
	GetRandomMeme(ctx context.Context, req *requester.Empty) (*requester.HighlightResponse, error)
}
