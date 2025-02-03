package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/chickey/blog/internal/models"
)

// commentCreator represents a type capable of reading a comment from storage and
// returning it or an error.
type commentCreator interface {
	CreateComment(ctx context.Context, comment models.Comment) (models.Comment, error)
}

// @Summary		Create Comment
// @Description	Creates a Comment
// @Tags			comment
// @Accept			json
// @Produce		json
// @Param			request	body		CommentRequest	true	"Comment to Create"
// @Success		200		{object}	uint
// @Failure		400		{object}	string
// @Failure		404		{object}	string
// @Failure		500		{object}	string
// @Router			/comment  [POST]
func HandleCreateComment(logger *slog.Logger, commentCreator commentCreator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Request validation
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
		}

		modelRequest := models.Comment{
			UserID:  request.UserID,
			BlogID:  request.BlogID,
			Message: request.Message,
		}
		// Read the comment
		comment, err := commentCreator.CreateComment(ctx, modelRequest)
		if err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to create comment",
				slog.String("error", err.Error()),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Convert our models.Comment domain model into a response model.
		response := CommentResponse{
			UserID:      comment.UserID,
			BlogID:      comment.BlogID,
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
