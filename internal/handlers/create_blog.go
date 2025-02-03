package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/chickey/blog/internal/models"
)

// blogCreator represents a type capable of reading a blog from storage and
// returning it or an error.
type blogCreator interface {
	CreateBlog(ctx context.Context, blog models.Blog) (models.Blog, error)
}

//	@Summary		Create Blog
//	@Description	Creates a Blog
//	@Tags			blog
//	@Accept			json
//	@Produce		json
//	@Param			request	body		BlogRequest	true	"Blog to Create"
//	@Success		200		{object}	uint
//	@Failure		400		{object}	string
//	@Failure		404		{object}	string
//	@Failure		500		{object}	string
//	@Router			/blog  [POST]
func HandleCreateBlog(logger *slog.Logger, blogCreator blogCreator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Request validation
		request, problems, err := decodeValid[*BlogRequest](r)

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
				slog.String("Validation failures:", fmt.Sprintf("%v", problems)),
			)
		}

		modelRequest := models.Blog{
			AuthorID: request.AuthorID,
			Title:    request.Title,
			Score:    request.Score,
		}
		// Read the blog
		blog, err := blogCreator.CreateBlog(ctx, modelRequest)
		if err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to create blog",
				slog.String("error", err.Error()),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Convert our models.Blog domain model into a response model.
		response := BlogResponse{
			ID:          blog.ID,
			AuthorID:    blog.AuthorID,
			Title:       blog.Title,
			Score:       blog.Score,
			CreatedDate: blog.CreatedDate,
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
