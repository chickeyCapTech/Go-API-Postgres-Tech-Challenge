package models

import "time"

type Blog struct {
	ID          uint
	AuthorID    uint
	Title       string
	Score       float32
	CreatedDate time.Time
}
