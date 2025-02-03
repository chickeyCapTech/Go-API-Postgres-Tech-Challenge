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

// blogUpdater represents a type capable of updating a blog and
// returning it or an error.
type blogUpdater interface {
	UpdateBlog(ctx context.Context, id uint64, patch models.Blog) (models.Blog, error)
}

//	@Summary		Update Blog
//	@Description	Update Blog by ID
//	@Tags			blog
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"Blog ID"
//	@Param			request	body		BlogRequest	true	"Blog to Create"
//	@Success		200		{object}	models.Blog
//	@Failure		400		{object}	string
//	@Failure		404		{object}	string
//	@Failure		500		{object}	string
//	@Router			/blog/{id}  [PUT]
func HandleUpdateBlog(logger *slog.Logger, blogUpdater blogUpdater) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Read id from path parameters
		idStr := r.PathValue("id")

		// Convert the ID from string to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to parse id from url",
				slog.String("id", idStr),
				slog.String("error", err.Error()),
			)

			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// Read request body
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
				slog.String("Validation errors: ", fmt.Sprintf("%#v", problems)),
			)
		}

		modelRequest := models.Blog{
			AuthorID: request.AuthorID,
			Title:    request.Title,
			Score:    request.Score,
		}

		// Update the blog
		blog, err := blogUpdater.UpdateBlog(ctx, uint64(id), modelRequest)
		if err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to update blog",
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
