package handlers

import (
	"context"
	"unicode/utf8"
)

type CommentRequest struct {
	UserID  uint
	BlogID  uint
	Message string
}

func (r *CommentRequest) Valid(ctx context.Context) map[string]string {

	problems := make(map[string]string)

	if r.UserID == 0 {
		problems["UserID"] = "Invalid UserId"
	}
	if r.BlogID == 0 {
		problems["BlogID"] = "Invalid BlogId"
	}
	if r.Message == "" {
		problems["Message"] = "Message cannot be empty"
	}
	if utf8.RuneCountInString(r.Message) > 500 {
		problems["Message"] = "Message cannot greater than 500 characters"
	}

	return problems
}
