package entity

import "time"

type MessageEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
}

type Message struct {
	Timestamp time.Time
	Text      string
}
