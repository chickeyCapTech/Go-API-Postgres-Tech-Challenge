package services

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log/slog"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chickey/blog/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCommentsService_ListComments(t *testing.T) {
	testcases := map[string]struct {
		mockCalled     bool
		mockInputArgs  []driver.Value
		mockOutput     *sqlmock.Rows
		mockError      error
		input          uint64
		expectedOutput []models.Comment
		expectedError  error
	}{
		"happy path": {
			mockCalled:    true,
			mockInputArgs: []driver.Value{1},
			mockOutput: sqlmock.NewRows([]string{"user_id", "blog_id", "message", "created_date"}).
				AddRow(1, 1, "New Comment", testDate).
				AddRow(1, 2, "Good blog", testDate),
			mockError: nil,
			input:     1,
			expectedOutput: []models.Comment{
				{
					BlogID:      1,
					UserID:      1,
					Message:     "New Comment",
					CreatedDate: testDate,
				},
				{
					BlogID:      2,
					UserID:      1,
					Message:     "Good blog",
					CreatedDate: testDate,
				},
			},
			expectedError: nil,
		},
		"no rows returned": {
			mockCalled:     true,
			mockInputArgs:  []driver.Value{1},
			mockOutput:     sqlmock.NewRows([]string{}),
			mockError:      sql.ErrNoRows,
			input:          1,
			expectedOutput: []models.Comment{},
			expectedError: fmt.Errorf(
				"[in services.CommentsService.ListComment] failed to list comments: %w",
				sql.ErrNoRows,
			),
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			logger := slog.Default()

			if tc.mockCalled {
				mock.
					ExpectQuery(regexp.QuoteMeta(`
                        SELECT user_id, blog_id, message, created_date FROM comments
                    `)).
					WillReturnRows(tc.mockOutput).
					WillReturnError(tc.mockError)
			}

			commentService := NewCommentsService(logger, db)

			outputs, err := commentService.ListComments(context.TODO(), 0, 0)
			if !assert.ErrorIs(t, err, tc.expectedError) {
				t.Errorf("expected %v, got %v", tc.expectedError, err)
			}
			for i, output := range outputs {
				if output != tc.expectedOutput[i] {
					t.Errorf("expected %v, got %v", tc.expectedOutput[i], output)
				}
			}

			if tc.mockCalled {
				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			}
		})
	}
}

func TestCommentsService_CreateComment(t *testing.T) {
	testcases := map[string]struct {
		mockCalled     bool
		mockInputArgs  []driver.Value
		mockOutput     *sqlmock.Rows
		mockError      error
		input          models.Comment
		expectedOutput models.Comment
		expectedError  error
	}{
		"happy path": {
			mockCalled:    true,
			mockInputArgs: []driver.Value{1, 1, "Good blog"},
			mockOutput: sqlmock.NewRows([]string{"created_date"}).
				AddRow(testDate),
			mockError: nil,
			input: models.Comment{
				BlogID:  1,
				UserID:  1,
				Message: "Good blog",
			},
			expectedOutput: models.Comment{
				BlogID:      1,
				UserID:      1,
				Message:     "Good blog",
				CreatedDate: testDate,
			},
			expectedError: nil,
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			logger := slog.Default()

			if tc.mockCalled {
				mock.
					ExpectQuery(regexp.QuoteMeta(`
				   SELECT 1
					FROM users
					WHERE id = $1::int
				`)).
					WithArgs([]driver.Value{1}...).
					WillReturnRows(sqlmock.NewRows([]string{"?column?"}).
						AddRow(1)).
					WillReturnError(tc.mockError)
				mock.
					ExpectQuery(regexp.QuoteMeta(`
                       SELECT 1
						FROM blogs
						WHERE id = $1::int
                    `)).
					WithArgs([]driver.Value{1}...).
					WillReturnRows(sqlmock.NewRows([]string{"?column?"}).
						AddRow(1)).
					WillReturnError(tc.mockError)
				mock.
					ExpectQuery(regexp.QuoteMeta(`
                    	SELECT 1
						FROM comments
						WHERE blog_id = $1::int AND user_id = $2::int
                    `)).
					WithArgs([]driver.Value{1, 1}...).
					WillReturnRows(sqlmock.NewRows([]string{"?column?"}).
						AddRow(1)).
					WillReturnError(tc.mockError)

				mock.
					ExpectQuery(regexp.QuoteMeta(
						`INSERT INTO comments (user_id, blog_id, message) VALUES ($1, $2, $3) RETURNING created_date`)).
					WithArgs(tc.mockInputArgs...).
					WillReturnRows(tc.mockOutput).
					WillReturnError(tc.mockError)
			}

			commentService := NewCommentsService(logger, db)

			output, err := commentService.CreateComment(context.TODO(), tc.input)
			if err != tc.expectedError {
				t.Errorf("expected no error, got %v", err)
			}
			if output != tc.expectedOutput {
				t.Errorf("expected %v, got %v", tc.expectedOutput, output)
			}

			if tc.mockCalled {
				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			}
		})
	}
}

func TestCommentsService_UpdateComment(t *testing.T) {
	testcases := map[string]struct {
		mockCalled     bool
		mockInputArgs  []driver.Value
		mockOutput     *sqlmock.Rows
		mockError      error
		input          models.Comment
		expectedOutput models.Comment
		expectedError  error
	}{
		"happy path": {
			mockCalled:    true,
			mockInputArgs: []driver.Value{"Good blog", 1, 1},
			mockOutput: sqlmock.NewRows([]string{"created_date"}).
				AddRow(testDate),
			mockError: nil,
			input: models.Comment{
				BlogID:  1,
				UserID:  1,
				Message: "Good blog",
			},
			expectedOutput: models.Comment{
				BlogID:      1,
				UserID:      1,
				Message:     "Good blog",
				CreatedDate: testDate,
			},
			expectedError: nil,
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			logger := slog.Default()

			if tc.mockCalled {
				mock.
					ExpectQuery(regexp.QuoteMeta(`
				   SELECT 1
					FROM users
					WHERE id = $1::int
				`)).
					WithArgs([]driver.Value{1}...).
					WillReturnRows(sqlmock.NewRows([]string{"?column?"}).
						AddRow(1)).
					WillReturnError(tc.mockError)
				mock.
					ExpectQuery(regexp.QuoteMeta(`
                       SELECT 1
						FROM blogs
						WHERE id = $1::int
                    `)).
					WithArgs([]driver.Value{1}...).
					WillReturnRows(sqlmock.NewRows([]string{"?column?"}).
						AddRow(1)).
					WillReturnError(tc.mockError)

				mock.
					ExpectQuery(regexp.QuoteMeta(
						`UPDATE comments
						SET message = $1
						WHERE user_id = $2 AND blog_id = $3
						RETURNING created_date`)).
					WithArgs(tc.mockInputArgs...).
					WillReturnRows(tc.mockOutput).
					WillReturnError(tc.mockError)
			}

			commentService := NewCommentsService(logger, db)

			output, err := commentService.UpdateComment(context.TODO(), tc.input)
			if err != tc.expectedError {
				t.Errorf("expected no error, got %v", err)
			}
			if output != tc.expectedOutput {
				t.Errorf("expected %v, got %v", tc.expectedOutput, output)
			}

			if tc.mockCalled {
				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			}
		})
	}
}

func TestCommentsService_DeleteComment(t *testing.T) {
	testcases := map[string]struct {
		mockCalled     bool
		mockInputArgs  []driver.Value
		mockOutput     *sqlmock.Rows
		mockError      error
		input          uint64
		expectedOutput models.Comment
		expectedError  error
	}{
		"happy path": {
			mockCalled:    true,
			mockInputArgs: []driver.Value{1, 1},
			mockError:     nil,
			input:         1,
			expectedError: nil,
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			logger := slog.Default()

			if tc.mockCalled {

				mock.
					ExpectExec(regexp.QuoteMeta(`DELETE FROM comments WHERE user_id = $1::int AND blog_id = $2::int`)).
					WithArgs(tc.mockInputArgs...).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(tc.mockError)
			}

			commentService := NewCommentsService(logger, db)

			err = commentService.DeleteComment(context.TODO(), uint(tc.input), uint(tc.input))
			if err != tc.expectedError {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.mockCalled {
				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			}
		})
	}
}
