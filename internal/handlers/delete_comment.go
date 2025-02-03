package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
)

// uerDeleter represents a type capable of deleting a comment from storage
type commentDeleter interface {
	DeleteComment(ctx context.Context, user_id uint, blog_id uint) error
}

// @Summary		Delete Comment
// @Description	Delete Comment by ID
// @Tags			comment
// @Accept			json
// @Produce		json
// @Param			author_id	query	string	false	"Author Id"
// @Param			blog_id		query	string	false	"Blog Id"
// @Success		200
// @Failure		400	{object}	string
// @Failure		404	{object}	string
// @Failure		500	{object}	string
// @Router			/comment  [DELETE]
func HandleDeleteComment(logger *slog.Logger, commentDeleter commentDeleter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Validate query params
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

				http.Error(w, "Invalid User ID", http.StatusBadRequest)
				return
			}
		}

		// Read the comment
		err = commentDeleter.DeleteComment(ctx, uint(userId), uint(blogId))
		if err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to read comment",
				slog.String("error", err.Error()),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Encode the response model as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}
