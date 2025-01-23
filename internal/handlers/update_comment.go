package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/chickey/blog/internal/models"
)

// commentUpdater represents a type capable of updating a comment and
// returning it or an error.
type commentUpdater interface {
	UpdateComment(ctx context.Context, patch models.Comment) (models.Comment, error)
}

// @Summary		Update Comment
// @Description	Update Comment by ID
// @Tags			comment
// @Accept			json
// @Produce		json
// @Param			author_id	query		string			false	"Author Id"
// @Param			blog_id		query		string			false	"Blog Id"
// @Param			request		body		CommentRequest	true	"Blog to Create"
// @Success		200			{object}	models.Comment
// @Failure		400			{object}	string
// @Failure		404			{object}	string
// @Failure		500			{object}	string
// @Router			/comment  [PUT]
func HandleUpdateComment(logger *slog.Logger, commentUpdater commentUpdater) http.Handler {
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

				http.Error(w, "Invalid Blog ID", http.StatusBadRequest)
				return
			}
		}

		// Read request body
		request, problems, err := decodeValid[*CommentRequest](r)

		if err != nil && len(problems) == 0 {
			logger.ErrorContext(
				r.Context(),
				"failed to decode request",
				slog.String("error", err.Error()))

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		if len(problems) > 0 {
			logger.ErrorContext(
				r.Context(),
				"Validation error",
				slog.String("Validation errors: ", fmt.Sprintf("%#v", problems)),
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		// Validae query params and body matches
		if request.BlogID != uint(blogId) || request.UserID != uint(userId) {
			logger.ErrorContext(
				r.Context(),
				"Validation error: Query Param does not match Request Body",
			)
			http.Error(w, "Invalid Request", http.StatusBadRequest)
		}

		modelRequest := models.Comment{
			UserID:  request.UserID,
			BlogID:  request.BlogID,
			Message: request.Message,
		}

		// Update the comment
		comment, err := commentUpdater.UpdateComment(ctx, modelRequest)
		if err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to update comment",
				slog.String("error", err.Error()),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Convert our models.Comment domain model into a response model.
		response := CommentResponse{
			BlogID:      comment.BlogID,
			UserID:      comment.UserID,
			Message:     comment.Message,
			CreatedDate: comment.CreatedDate,
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
