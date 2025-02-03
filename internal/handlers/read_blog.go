package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/chickey/blog/internal/models"
)

// blogReader represents a type capable of reading a blog from storage and
// returning it or an error.
type blogReader interface {
	ReadBlog(ctx context.Context, id uint64) (models.Blog, error)
}

// @Summary		Read Blog
// @Description	Read Blog by ID
// @Tags			blog
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"Blog Id"
// @Success		200	{object}	models.Blog
// @Failure		400	{object}	string
// @Failure		404	{object}	string
// @Failure		500	{object}	string
// @Router			/blog/{id}  [GET]
func HandleReadBlog(logger *slog.Logger, blogReader blogReader) http.Handler {
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

		// Read the blog
		blog, err := blogReader.ReadBlog(ctx, uint64(id))
		if err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to read blog",
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
