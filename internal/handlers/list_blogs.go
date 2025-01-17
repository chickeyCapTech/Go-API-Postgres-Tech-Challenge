package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/chickey/blog/internal/models"
)

// blogReader represents a type capable of reading a blog from storage and
// returning it or an error.
type blogsReader interface {
	ListBlogs(ctx context.Context, title string) ([]models.Blog, error)
}

// listBlogsResponse represents the response for listing blogs.
type listBlogsResponse struct {
	Blogs []readBlogResponse
}

// @Summary		List Blogs
// @Description	List All Blogs
// @Tags			blogs
// @Accept			json
// @Produce		json
// @Param title query string false "query param"
// @Success		200	{array}		models.Blog
// @Failure		400	{object}	string
// @Failure		404	{object}	string
// @Failure		500	{object}	string
// @Router			/blogs  [GET]
func HandleListBlogs(logger *slog.Logger, blogsReader blogsReader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		title := r.URL.Query().Get("title")

		// Read the blog
		blogs, err := blogsReader.ListBlogs(ctx, title)
		if err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to list blogs",
				slog.String("error", err.Error()),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Convert our models.Blog domain model into a response model.
		response := listBlogsResponse{
			Blogs: []readBlogResponse{},
		}

		for _, blog := range blogs {
			newBlog := readBlogResponse{
				ID:          blog.ID,
				AuthorID:    blog.AuthorID,
				Title:       blog.Title,
				Score:       blog.Score,
				CreatedDate: blog.CreatedDate,
			}
			response.Blogs = append(response.Blogs, newBlog)
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
