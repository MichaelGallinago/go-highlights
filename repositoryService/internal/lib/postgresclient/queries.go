package postgresclient

import (
	"context"
	"github.com/jackc/pgx/v5"
	"repositoryService/internal/core/entity"
)

// InsertMeme добавляет мем в БД
func (db *PostgresClient) InsertMeme(ctx context.Context, timestamp int64, text string) error {
	query := `
INSERT INTO memes (timestamp, text) 
VALUES ($1, $2) 
ON CONFLICT (text) 
DO UPDATE SET timestamp = EXCLUDED.timestamp
`
	_, err := db.Pool.Exec(ctx, query, timestamp, text)
	return err
}

// GetTopLongMemes возвращает топ мемов по длине текста
func (db *PostgresClient) GetTopLongMemes(ctx context.Context, limit int) ([]entity.Meme, error) {
	query := "SELECT timestamp, text FROM memes ORDER BY LENGTH(text) DESC LIMIT $1"
	rows, err := db.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memes, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Meme])
	if err != nil {
		return nil, err
	}
	return memes, nil
}

// SearchMemes ищет мемы по подстроке
func (db *PostgresClient) SearchMemes(ctx context.Context, query string) ([]entity.Meme, error) {
	sql := "SELECT timestamp, text FROM memes WHERE text ILIKE '%' || $1 || '%'"
	rows, err := db.Pool.Query(ctx, sql, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memes, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Meme])
	if err != nil {
		return nil, err
	}
	return memes, nil
}

// GetMemesByMonth возвращает мемы за указанный период (месяц)
func (db *PostgresClient) GetMemesByMonth(ctx context.Context, month int32) ([]entity.Meme, error) {
	query := `SELECT timestamp, text FROM memes WHERE EXTRACT(MONTH FROM TO_TIMESTAMP(timestamp)) = $1`
	rows, err := db.Pool.Query(ctx, query, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memes, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Meme])
	if err != nil {
		return nil, err
	}
	return memes, nil
}

// GetRandomMeme возвращает случайный мем
func (db *PostgresClient) GetRandomMeme(ctx context.Context) (entity.Meme, error) {
	var meme entity.Meme
	query := "SELECT timestamp, text FROM memes ORDER BY RANDOM() LIMIT 1"
	err := db.Pool.QueryRow(ctx, query).Scan(&meme.Timestamp, &meme.Text)
	if err != nil {
		return entity.Meme{}, err
	}
	return meme, nil
}
