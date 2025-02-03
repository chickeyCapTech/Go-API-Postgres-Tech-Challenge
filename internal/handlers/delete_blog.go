package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
)

// uerDeleter represents a type capable of deleting a blog from storage
type blogDeleter interface {
	DeleteBlog(ctx context.Context, id uint64) error
}

// @Summary		Delete Blog
// @Description	Delete Blog by ID
// @Tags			blog
// @Accept			json
// @Produce		json
// @Param			id	path	string	true	"Blog ID"
// @Success		200
// @Failure		400	{object}	string
// @Failure		404	{object}	string
// @Failure		500	{object}	string
// @Router			/blog/{id}  [DELETE]
func HandleDeleteBlog(logger *slog.Logger, blogDeleter blogDeleter) http.Handler {
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
		err = blogDeleter.DeleteBlog(ctx, uint64(id))
		if err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to read blog",
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
