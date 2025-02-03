package handlers

import (
	"context"
	"unicode/utf8"
)

// BlogRequest represents the request for creating a Blog.
type BlogRequest struct {
	AuthorID uint    `json:"authorid"`
	Title    string  `json:"title"`
	Score    float32 `json:"score"`
}

func (r *BlogRequest) Valid(ctx context.Context) map[string]string {

	problems := make(map[string]string)

	if r.AuthorID == 0 {
		problems["AuthorId"] = "Invalid AuthorId"
	}
	if r.Score < 0 || r.Score > 10 {
		problems["Score"] = "Score cannot be less than 0 or greater than 10"
	}
	if r.Title == "" {
		problems["Title"] = "Title cannot be empty"
	}
	if utf8.RuneCountInString(r.Title) > 100 {
		problems["Title"] = "Title cannot be greater than 100 characters"
	}

	return problems
}
