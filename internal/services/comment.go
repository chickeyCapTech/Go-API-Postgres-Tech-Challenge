package services

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/chickey/blog/internal/models"
)

// CommentsService is a service capable of performing CRUD operations for
// models.Comment models.
type CommentsService struct {
	logger *slog.Logger
	db     *sql.DB
}

// NewCommentsService creates a new CommentsService and returns a pointer to it.
func NewCommentsService(logger *slog.Logger, db *sql.DB) *CommentsService {
	return &CommentsService{
		logger: logger,
		db:     db,
	}
}

// CreateComment attempts to create the provided comment, returning a fully hydrated
// models.Comment or an error.
func (s *CommentsService) CreateComment(ctx context.Context, comment models.Comment) (models.Comment, error) {
	s.logger.DebugContext(ctx, "Creating comment", "Blog ID", comment.BlogID, "UserId", comment.UserID)

	//validate user_id exists in user table
	author := s.db.QueryRowContext(
		ctx,
		`
			SELECT 1
			FROM users
			WHERE id = $1::int
			`,
		comment.UserID,
	)

	var authorExists int

	err := author.Scan(&authorExists)

	if err != nil {
		return models.Comment{}, fmt.Errorf(
			"[in services.CommentsService.UpdateComment] failed to create comment: %w",
			err,
		)
	}

	//validate blog exists with blog_id
	blog := s.db.QueryRowContext(
		ctx,
		`
			SELECT 1
			FROM blogs
			WHERE id = $1::int
			`,
		comment.BlogID,
	)

	var blogExists int

	err = blog.Scan(&blogExists)

	if err != nil {
		return models.Comment{}, fmt.Errorf(
			"[in services.CommentsService.UpdateComment] failed to create comment: %w",
			err,
		)
	}

	//validate comment does not already exist with blog_id and user_id
	commentCheck := s.db.QueryRowContext(
		ctx,
		`
			SELECT 1
			FROM comments
			WHERE blog_id = $1::int AND user_id = $2::int
			`,
		comment.BlogID,
		comment.UserID,
	)

	var commentExists int

	err = commentCheck.Scan(&commentExists)

	if err != nil {
		return models.Comment{}, fmt.Errorf(
			"[in services.CommentsService.UpdateComment] failed to create comment: %w",
			err,
		)
	}

	// Create new comment entry in comment table
	result := s.db.QueryRowContext(
		ctx,
		`
		INSERT INTO comments (user_id, blog_id, message) VALUES ($1, $2, $3) RETURNING created_date
		`,
		comment.UserID,
		comment.BlogID,
		comment.Message,
	)

	err = result.Scan(&comment.CreatedDate)

	if err != nil {
		return models.Comment{}, fmt.Errorf(
			"[in services.CommentsService.CreateComment] failed to create comment: %w",
			err,
		)
	}

	return comment, nil
}

// UpdateComment attempts to perform an update of the comment with the provided id,
// updating, it to reflect the properties on the provided patch object. A
// models.Comment or an error.
func (s *CommentsService) UpdateComment(ctx context.Context, patch models.Comment) (models.Comment, error) {
	s.logger.DebugContext(ctx, "Updating comment", "Blog ID", patch.BlogID, "UserId", patch.UserID)

	//validate user_id exists in user table
	author := s.db.QueryRowContext(
		ctx,
		`
		SELECT 1
		FROM users
		WHERE id = $1::int
        `,
		patch.UserID,
	)

	var authorExists int

	err := author.Scan(&authorExists)

	if err != nil {
		return models.Comment{}, fmt.Errorf(
			"[in services.CommentsService.UpdateComment] failed to update comment: %w",
			err,
		)
	}

	//validate blog exists with blog_id
	blog := s.db.QueryRowContext(
		ctx,
		`
		SELECT 1
		FROM blogs
		WHERE id = $1::int
        `,
		patch.BlogID,
	)

	var blogExists int

	err = blog.Scan(&blogExists)

	if err != nil {
		return models.Comment{}, fmt.Errorf(
			"[in services.CommentsService.UpdateComment] failed to update comment: %w",
			err,
		)
	}

	// should we be able to update created date?
	row := s.db.QueryRowContext(
		ctx,
		`
		UPDATE comments
		SET message = $1
		WHERE user_id = $2 AND blog_id = $3
		RETURNING created_date
		`,
		patch.Message,
		patch.UserID,
		patch.BlogID,
	)
	err = row.Scan(&patch.CreatedDate)

	if err != nil {
		return models.Comment{}, fmt.Errorf(
			"[in services.CommentsService.UpdateComment] failed to update blog: %w",
			err,
		)
	}

	return patch, nil
}

// DeleteComment attempts to delete the comment with the provided id. An error is
// returned if the delete fails.
func (s *CommentsService) DeleteComment(ctx context.Context, userId uint, blogId uint) error {
	s.logger.DebugContext(ctx, "Deleteing comment", "User Id", userId, "Blog Id", blogId)

	//DELETE from comment from comments
	_, err := s.db.ExecContext(
		ctx,
		`
		DELETE FROM comments WHERE user_id = $1::int AND blog_id = $2::int
		`,
		userId,
		blogId,
	)

	if err != nil {
		return fmt.Errorf(
			"[in services.CommentsService.DeleteComment] failed to delete comment: %w",
			err,
		)
	}

	return nil
}

// ListComments attempts to list all comments in the database. A slice of models.Comment
// or an error is returned.
// WHY DO WE NEED ID?
func (s *CommentsService) ListComments(ctx context.Context, userId uint, blogId uint) ([]models.Comment, error) {
	s.logger.DebugContext(ctx, "Listing comments")

	//Build query based on query params
	baseQuery := "SELECT user_id, blog_id, message, created_date FROM comments"

	conditions := []string{}
	args := []any{}
	i := 1

	if userId > 0 {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", i))
		args = append(args, userId)
		i++
	}
	if blogId > 0 {
		conditions = append(conditions, fmt.Sprintf("blog_id = $%d", i))
		args = append(args, blogId)
		i++
	}

	if len(conditions) > 0 {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(conditions, " AND "))
	}

	rows, err := s.db.QueryContext(
		ctx,
		baseQuery,
		args...,
	)

	if err != nil {
		return []models.Comment{}, fmt.Errorf(
			"[in services.CommentsService.ListComment] failed to read comments: %w",
			err,
		)
	}

	var comments []models.Comment

	// error handling for error in for loop and then for empty error out of for loop?
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.UserID, &comment.BlogID, &comment.Message, &comment.CreatedDate)
		if err != nil {
			return []models.Comment{}, fmt.Errorf(
				"[in services.CommentsService.ListComment] failed to read comments: %w",
				err,
			)
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return []models.Comment{}, fmt.Errorf(
			"[in services.CommentsService.ListComment] failed to read comments: %w",
			err,
		)
	}

	return comments, nil
}
