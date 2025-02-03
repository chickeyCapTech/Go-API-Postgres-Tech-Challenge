package handlers

import (
	"time"
)

// BlogResponse represents the response for creating a Blog.
type BlogResponse struct {
	ID          uint      `json:"id"`
	AuthorID    uint      `json:"authorid"`
	Title       string    `json:"title"`
	Score       float32   `json:"score"`
	CreatedDate time.Time `json:"createddate"`
}
