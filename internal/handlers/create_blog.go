package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/chickey/blog/internal/models"
)

// blogCreator represents a type capable of reading a blog from storage and
// returning it or an error.
type blogCreator interface {
	CreateBlog(ctx context.Context, blog models.Blog) (models.Blog, error)
}

// createBlogRequest represents the request for creating a Blog.
type createBlogRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// createBlogResponse represents the response for creating a Blog.
type createBlogResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary		Create Blog
// @Description	Creates a Blog
// @Tags			blogs
// @Accept			json
// @Produce		json
// @Param			request	body		createBlogRequest	true	"Blog to Create"
// @Success		200		{object}	uint
// @Failure		400		{object}	string
// @Failure		404		{object}	string
// @Failure		500		{object}	string
// @Router			/blogs  [POST]
func HandleCreateBlog(logger *slog.Logger, blogCreator blogCreator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var request models.Blog
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to decode request",
				slog.String("error", err.Error()))

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		// Read the blog
		blog, err := blogCreator.CreateBlog(ctx, request)
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
		response := createBlogResponse{
			ID: blog.ID,
			//TODO
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
