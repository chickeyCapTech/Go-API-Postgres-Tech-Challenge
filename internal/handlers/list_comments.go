package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/chickey/blog/internal/models"
)

// commentReader represents a type capable of reading a comment from storage and
// returning it or an error.
type commentsLister interface {
	ListComments(ctx context.Context, authorId uint, blogId uint) ([]models.Comment, error)
}

// listCommentsResponse represents the response for listing comments.
type listCommentsResponse struct {
	Comments []CommentResponse
}

// @Summary		List Comments
// @Description	List All Comments
// @Tags			comment
// @Accept			json
// @Produce		json
// @Param			author_id	query		string	false	"Author Id"
// @Param			blog_id		query		string	false	"Blog Id"
// @Success		200			{array}		models.Comment
// @Failure		400			{object}	string
// @Failure		404			{object}	string
// @Failure		500			{object}	string
// @Router			/comment  [GET]
func HandleListComments(logger *slog.Logger, commentsLister commentsLister) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userIdStr := r.URL.Query().Get("author_id")
		blogIdStr := r.URL.Query().Get("blog_id")

		var userId int
		var err error

		if userIdStr != "" {
			userId, err = strconv.Atoi(userIdStr)
			if err != nil {
				logger.ErrorContext(
					r.Context(),
					"failed to get valid User id from query param",
					slog.String("id", userIdStr),
					slog.String("error", err.Error()),
				)

				http.Error(w, "Invalid User ID", http.StatusBadRequest)
				return
			}
		}

		var blogId int

		if blogIdStr != "" {
			blogId, err = strconv.Atoi(blogIdStr)
			if err != nil {
				logger.ErrorContext(
					r.Context(),
					"failed to get valid Blog Id from query param",
					slog.String("id", blogIdStr),
					slog.String("error", err.Error()),
				)

				http.Error(w, "Invalid Blog ID", http.StatusBadRequest)
				return
			}
		}

		// Read the comment
		comments, err := commentsLister.ListComments(ctx, uint(userId), uint(blogId))
		if err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to list comments",
				slog.String("error", err.Error()),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Convert our models.Comment domain model into a response model.
		response := listCommentsResponse{
			Comments: []CommentResponse{},
		}

		for _, comment := range comments {
			newComment := CommentResponse{
				BlogID:      comment.BlogID,
				UserID:      comment.UserID,
				Message:     comment.Message,
				CreatedDate: comment.CreatedDate,
			}
			response.Comments = append(response.Comments, newComment)
		}

		// Encode the response model as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to encode response",
				slog.String("error", err.Error()))

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})
}
