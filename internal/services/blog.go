package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/chickey/blog/internal/models"
)

// BlogsService is a service capable of performing CRUD operations for
// models.Blog models.
type BlogsService struct {
	logger *slog.Logger
	db     *sql.DB
}

// NewBlogsService creates a new BlogsService and returns a pointer to it.
func NewBlogsService(logger *slog.Logger, db *sql.DB) *BlogsService {
	return &BlogsService{
		logger: logger,
		db:     db,
	}
}

// // CreateBlog attempts to create the provided blog, returning a fully hydrated
// // models.Blog or an error.
// func (s *BlogsService) CreateBlog(ctx context.Context, blog models.Blog) (models.Blog, error) {
// 	s.logger.DebugContext(ctx, "Creating blog", "name", blog.Name)

// 	result, err := s.db.ExecContext(
// 		ctx,
// 		`
// 		INSERT INTO blogs (name, email, password) VALUE ($1::string, $2::string, $3::string)
// 		`,
// 		blog.Name,
// 		blog.Email,
// 		blog.Password,
// 	)

// 	if err != nil {
// 		return models.Blog{}, fmt.Errorf(
// 			"[in services.BlogsService.CreateBlog] failed to create blog: %w",
// 			err,
// 		)
// 	}

// 	id, err := result.LastInsertId()

// 	if err != nil {
// 		return models.Blog{}, fmt.Errorf(
// 			"[in services.BlogsService.CreateBlog] failed to create blog: %w",
// 			err,
// 		)
// 	}

// 	blog.ID = uint(id)

// 	return blog, nil
// }

// ReadBlog attempts to read a blog from the database using the provided id. A
// fully hydrated models.Blog or error is returned.
func (s *BlogsService) ReadBlog(ctx context.Context, id uint64) (models.Blog, error) {
	s.logger.DebugContext(ctx, "Reading blog", "id", id)

	row := s.db.QueryRowContext(
		ctx,
		`
		SELECT id,
		       author_id,
		       title,
		       score,
			   created_date
		FROM blogs
		WHERE id = $1::int
        `,
		id,
	)

	var blog models.Blog

	err := row.Scan(&blog.ID, &blog.AuthorID, &blog.Title, &blog.Score, &blog.CreatedDate)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Blog{}, nil
		default:
			return models.Blog{}, fmt.Errorf(
				"[in services.BlogsService.ReadBlog] failed to read blog: %w",
				err,
			)
		}
	}

	return blog, nil
}

// // UpdateBlog attempts to perform an update of the blog with the provided id,
// // updating, it to reflect the properties on the provided patch object. A
// // models.Blog or an error.
// func (s *BlogsService) UpdateBlog(ctx context.Context, id uint64, patch models.Blog) (models.Blog, error) {
// 	s.logger.DebugContext(ctx, "Updating blog", "id", id)

// 	_, err := s.db.ExecContext(
// 		ctx,
// 		`
// 		UPDATE blogs
// 		SET name = $1, email = $2, password = $3
// 		WHERE id = $4
// 		`,
// 		patch.Name,
// 		patch.Email,
// 		patch.Password,
// 		id,
// 	)

// 	if err != nil {
// 		return models.Blog{}, fmt.Errorf(
// 			"[in services.BlogsService.UpdateBlog] failed to update blog: %w",
// 			err,
// 		)
// 	}

// 	return patch, nil
// }

// // DeleteBlog attempts to delete the blog with the provided id. An error is
// // returned if the delete fails.
// func (s *BlogsService) DeleteBlog(ctx context.Context, id uint64) error {
// 	s.logger.DebugContext(ctx, "Creating blog", "id", id)

// 	_, err := s.db.ExecContext(
// 		ctx,
// 		`
// 		DELETE FROM blogs WHERE id = $1::int
// 		`,
// 		id,
// 	)

// 	if err != nil {
// 		return fmt.Errorf(
// 			"[in services.BlogsService.DeleteBlog] failed to delete blog: %w",
// 			err,
// 		)
// 	}

// 	return nil
// }

// ListBlogs attempts to list all blogs in the database. A slice of models.Blog
// or an error is returned.
// WHY DO WE NEED ID?
func (s *BlogsService) ListBlogs(ctx context.Context, title string) ([]models.Blog, error) {
	s.logger.DebugContext(ctx, "Listing blogs")

	rows, err := s.db.QueryContext(
		ctx,
		`
		SELECT id,
		       name,
		       email,
		       password
		FROM blogs
        `,
	)

	var blogs []models.Blog

	// error handling for error in for loop and then for empty error out of for loop?
	// NEED TO HANDLE ERROR FOR EMPTY RESULTS
	for rows.Next() {
		var blog models.Blog
		err := rows.Scan(&blog.ID, &blog.AuthorID, &blog.Title, &blog.Score, &blog.CreatedDate)
		if err != nil {
			return []models.Blog{}, fmt.Errorf(
				"[in services.BlogsService.ListBlog] failed to read blogs: %w",
				err,
			)
		}
		if title == "" || blog.Title == title {
			blogs = append(blogs, blog)
		}

	}

	if err = rows.Err(); err != nil {
		return []models.Blog{}, fmt.Errorf(
			"[in services.BlogsService.ListBlog] failed to read blogs: %w",
			err,
		)
	}

	return blogs, nil
}
