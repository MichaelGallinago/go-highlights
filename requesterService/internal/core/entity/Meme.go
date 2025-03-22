package entity

import "time"

type Meme struct {
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
}
