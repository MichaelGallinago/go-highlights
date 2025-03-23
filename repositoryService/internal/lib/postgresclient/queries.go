package postgresclient

import (
	"context"
	"repositoryService/internal/core/entity"
	"time"
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

	var memes []entity.Meme
	for rows.Next() {
		var meme entity.Meme
		var timestampInt int64
		if err := rows.Scan(&timestampInt, &meme.Text); err != nil {
			return nil, err
		}

		meme.Timestamp = toTime(timestampInt)
		memes = append(memes, meme)
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

	var memes []entity.Meme
	for rows.Next() {
		var meme entity.Meme
		var timestampInt int64
		if err := rows.Scan(&timestampInt, &meme.Text); err != nil {
			return nil, err
		}

		meme.Timestamp = toTime(timestampInt)
		memes = append(memes, meme)
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

	var memes []entity.Meme
	for rows.Next() {
		var meme entity.Meme
		var timestampInt int64
		if err := rows.Scan(&timestampInt, &meme.Text); err != nil {
			return nil, err
		}

		meme.Timestamp = toTime(timestampInt)
		memes = append(memes, meme)
	}
	return memes, nil
}

// GetRandomMeme возвращает случайный мем
func (db *PostgresClient) GetRandomMeme(ctx context.Context) (entity.Meme, error) {
	query := "SELECT timestamp, text FROM memes ORDER BY RANDOM() LIMIT 1"
	row := db.Pool.QueryRow(ctx, query)

	var meme entity.Meme
	var timestampInt int64
	err := row.Scan(&timestampInt, &meme.Text)
	if err != nil {
		return entity.Meme{}, err
	}
	meme.Timestamp = toTime(timestampInt)

	return meme, nil
}

func toTime(timestampInt int64) time.Time {
	return time.Unix(timestampInt, 0)
}
