package handlers

import "time"

type CommentResponse struct {
	UserID      uint
	BlogID      uint
	Message     string
	CreatedDate time.Time
}
