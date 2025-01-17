package models

import "time"

type Comment struct {
	UserID      uint
	BlogID      uint
	Message     string
	CreatedDate time.Time
}
