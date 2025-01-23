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

// CreateBlog attempts to create the provided blog, returning a fully hydrated
// models.Blog or an error.
func (s *BlogsService) CreateBlog(ctx context.Context, blog models.Blog) (models.Blog, error) {
	s.logger.DebugContext(ctx, "Creating blog", "name", blog.Title)

	//validate authod_id exists in user table
	author := s.db.QueryRowContext(
		ctx,
		`
			SELECT id,
				   name,
				   email,
				   password
			FROM users
			WHERE id = $1::int
			`,
		blog.AuthorID,
	)

	var user models.User

	err := author.Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		return models.Blog{}, fmt.Errorf(
			"[in services.BlogsService.UpdateBlog] failed to update blog: %w",
			err,
		)
	}

	// Create new blog entry in blog table
	result := s.db.QueryRowContext(
		ctx,
		`
		INSERT INTO blogs (author_id, title, score) VALUES ($1, $2, $3) RETURNING id, created_date
		`,
		blog.AuthorID,
		blog.Title,
		blog.Score,
	)

	err = result.Scan(&blog.ID, &blog.CreatedDate)

	if err != nil {
		return models.Blog{}, fmt.Errorf(
			"[in services.BlogsService.CreateBlog] failed to create blog: %w",
			err,
		)
	}

	return blog, nil
}

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

// UpdateBlog attempts to perform an update of the blog with the provided id,
// updating, it to reflect the properties on the provided patch object. A
// models.Blog or an error.
func (s *BlogsService) UpdateBlog(ctx context.Context, id uint64, patch models.Blog) (models.Blog, error) {
	s.logger.DebugContext(ctx, "Updating blog", "id", id)

	//validate authod_id exists in user table
	author := s.db.QueryRowContext(
		ctx,
		`
		SELECT 1
		FROM users
		WHERE id = $1::int
        `,
		id,
	)

	var exists int

	err := author.Scan(&exists)

	if err != nil {
		return models.Blog{}, fmt.Errorf(
			"[in services.BlogsService.UpdateBlog] failed to update blog: %w",
			err,
		)
	}

	// should we be able to update created date?
	row := s.db.QueryRowContext(
		ctx,
		`
		UPDATE blogs
		SET author_id = $1, title = $2, score = $3
		WHERE id = $4
		RETURNING created_date
		`,
		patch.AuthorID,
		patch.Title,
		patch.Score,
		id,
	)
	err = row.Scan(&patch.CreatedDate)

	if err != nil {
		return models.Blog{}, fmt.Errorf(
			"[in services.BlogsService.UpdateBlog] failed to update blog: %w",
			err,
		)
	}
	patch.ID = uint(id)
	return patch, nil
}

// DeleteBlog attempts to delete the blog with the provided id. An error is
// returned if the delete fails.
func (s *BlogsService) DeleteBlog(ctx context.Context, id uint64) error {
	s.logger.DebugContext(ctx, "Deleting blog", "id", id)

	//DB transaction

	//DELETE from blog from blogs
	_, err := s.db.ExecContext(
		ctx,
		`
		DELETE FROM blogs WHERE id = $1::int
		`,
		id,
	)

	if err != nil {
		return fmt.Errorf(
			"[in services.BlogsService.DeleteBlog] failed to delete blog: %w",
			err,
		)
	}

	//DELETE all comments with blog_id of deleted blog
	_, err = s.db.ExecContext(
		ctx,
		`
		DELETE FROM comments WHERE blog_id = $1::int
		`,
		id,
	)

	if err != nil {
		return fmt.Errorf(
			"[in services.BlogsService.DeleteBlog] failed to delete comments of deleted blog: %w",
			err,
		)
	}

	return nil
}

// ListBlogs attempts to list all blogs in the database. A slice of models.Blog
// or an error is returned.
// WHY DO WE NEED ID?
func (s *BlogsService) ListBlogs(ctx context.Context, title string) ([]models.Blog, error) {
	s.logger.DebugContext(ctx, "Listing blogs")

	rows, err := s.db.QueryContext(
		ctx,
		`
		SELECT id, author_id, title, score, created_date
		FROM blogs
        `,
	)

	if err != nil {
		return []models.Blog{}, fmt.Errorf(
			"[in services.CommentsService.ListComment] failed to read blogs: %w",
			err,
		)
	}

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
		if len(title) == 0 || blog.Title == title {
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
