package services

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log/slog"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chickey/blog/internal/models"
)

var testDate time.Time = time.Date(2025, 1, 21, 11, 12, 11, 11, time.UTC)

func TestBlogsService_ReadBlog(t *testing.T) {
	testcases := map[string]struct {
		mockCalled     bool
		mockInputArgs  []driver.Value
		mockOutput     *sqlmock.Rows
		mockError      error
		input          uint64
		expectedOutput models.Blog
		expectedError  error
	}{
		"happy path": {
			mockCalled:    true,
			mockInputArgs: []driver.Value{1},
			mockOutput: sqlmock.NewRows([]string{"id", "author_id", "title", "score", "created_date"}).
				AddRow(1, 1, "Book Title", 8.2, testDate),
			mockError: nil,
			input:     1,
			expectedOutput: models.Blog{
				ID:          1,
				AuthorID:    1,
				Title:       "Book Title",
				Score:       8.2,
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
                        SELECT id,
		       author_id,
		       title,
		       score,
			   created_date
		FROM blogs
		WHERE id = $1::int
                    `)).
					WithArgs(tc.mockInputArgs...).
					WillReturnRows(tc.mockOutput).
					WillReturnError(tc.mockError)
			}

			blogService := NewBlogsService(logger, db)

			output, err := blogService.ReadBlog(context.TODO(), tc.input)
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
func TestBlogsService_ListBlogs(t *testing.T) {
	testcases := map[string]struct {
		mockCalled     bool
		mockInputArgs  []driver.Value
		mockOutput     *sqlmock.Rows
		mockError      error
		input          uint64
		expectedOutput []models.Blog
		expectedError  error
	}{
		"happy path": {
			mockCalled:    true,
			mockInputArgs: []driver.Value{1},
			mockOutput: sqlmock.NewRows([]string{"id", "author_id", "title", "score", "created_date"}).
				AddRow(1, 1, "Book Title", 8.2, testDate).
				AddRow(2, 1, "New Book", 7.4, testDate),
			mockError: nil,
			input:     1,
			expectedOutput: []models.Blog{
				{
					ID:          1,
					AuthorID:    1,
					Title:       "Book Title",
					Score:       8.2,
					CreatedDate: testDate,
				},
				{
					ID:          2,
					AuthorID:    1,
					Title:       "New Book",
					Score:       7.4,
					CreatedDate: testDate,
				},
			},
			expectedError: nil,
		},
		"title filter": {
			mockCalled:    true,
			mockInputArgs: []driver.Value{1},
			mockOutput: sqlmock.NewRows([]string{"id", "author_id", "title", "score", "created_date"}).
				AddRow(1, 1, "Book Title", 8.2, testDate).
				AddRow(2, 1, "New Book", 7.4, testDate),
			mockError: nil,
			input:     1,
			expectedOutput: []models.Blog{
				{
					ID:          1,
					AuthorID:    1,
					Title:       "Book Title",
					Score:       8.2,
					CreatedDate: testDate,
				},
			},
			expectedError: nil,
		},
		"no rows returned": {
			mockCalled:    true,
			mockInputArgs: []driver.Value{1},
			mockOutput: sqlmock.NewRows([]string{"id", "author_id", "title", "score", "created_date"}).
				AddRow(1, 1, "Book Title", 8.2, testDate).
				AddRow(2, 1, "New Book", 7.4, testDate),
			mockError:      sql.ErrNoRows,
			input:          1,
			expectedOutput: []models.Blog{},
			expectedError: fmt.Errorf(
				"[in services.CommentsService.ListComment] failed to list blogs: %w",
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
                        SELECT id, author_id, title, score, created_date
					FROM blogs
                    `)).
					WillReturnRows(tc.mockOutput).
					WillReturnError(tc.mockError)
			}

			blogService := NewBlogsService(logger, db)

			outputs, err := blogService.ListBlogs(context.TODO(), "Book Title")
			if err != tc.expectedError {
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

func TestBlogsService_CreateBlog(t *testing.T) {
	testcases := map[string]struct {
		mockCalled     bool
		mockInputArgs  []driver.Value
		mockOutput     *sqlmock.Rows
		mockError      error
		input          models.Blog
		expectedOutput models.Blog
		expectedError  error
	}{
		"happy path": {
			mockCalled:    true,
			mockInputArgs: []driver.Value{1, "Book Title", float32(8.2)},
			mockOutput: sqlmock.NewRows([]string{"id", "created_date"}).
				AddRow(1, testDate),
			mockError: nil,
			input: models.Blog{
				AuthorID: 1,
				Title:    "Book Title",
				Score:    8.2,
			},
			expectedOutput: models.Blog{
				ID:          1,
				AuthorID:    1,
				Title:       "Book Title",
				Score:       8.2,
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
                       SELECT id,
                               name,
                               email,
                               password
                        FROM users
                        WHERE id = $1::int
                    `)).
					WithArgs([]driver.Value{1}...).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
						AddRow(1, "john", "john@me.com", "password123!")).
					WillReturnError(tc.mockError)

				mock.
					ExpectQuery(regexp.QuoteMeta(
						`INSERT INTO blogs (author_id, title, score) VALUES ($1, $2, $3) RETURNING id, created_date`)).
					WithArgs(tc.mockInputArgs...).
					WillReturnRows(tc.mockOutput).
					WillReturnError(tc.mockError)
			}

			blogService := NewBlogsService(logger, db)

			output, err := blogService.CreateBlog(context.TODO(), tc.input)
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

func TestBlogsService_UpdateBlog(t *testing.T) {
	testcases := map[string]struct {
		mockCalled     bool
		mockInputArgs  []driver.Value
		mockOutput     *sqlmock.Rows
		mockError      error
		input          models.Blog
		expectedOutput models.Blog
		expectedError  error
	}{
		"happy path": {
			mockCalled:    true,
			mockInputArgs: []driver.Value{1, "Book Title", float32(8.2), 1},
			mockOutput: sqlmock.NewRows([]string{"created_date"}).
				AddRow(testDate),
			mockError: nil,
			input: models.Blog{
				AuthorID: 1,
				Title:    "Book Title",
				Score:    8.2,
			},
			expectedOutput: models.Blog{
				ID:          1,
				AuthorID:    1,
				Title:       "Book Title",
				Score:       8.2,
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
					ExpectQuery(regexp.QuoteMeta(
						`UPDATE blogs
		SET author_id = $1, title = $2, score = $3
		WHERE id = $4
		RETURNING created_date`)).
					WithArgs(tc.mockInputArgs...).
					WillReturnRows(tc.mockOutput).
					WillReturnError(tc.mockError)
			}

			blogService := NewBlogsService(logger, db)

			output, err := blogService.UpdateBlog(context.TODO(), 1, tc.input)
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

func TestBlogsService_DeleteBlog(t *testing.T) {
	testcases := map[string]struct {
		mockCalled     bool
		mockInputArgs  []driver.Value
		mockOutput     *sqlmock.Rows
		mockError      error
		input          uint64
		expectedOutput models.Blog
		expectedError  error
	}{
		"happy path": {
			mockCalled:    true,
			mockInputArgs: []driver.Value{1},
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
					ExpectExec(regexp.QuoteMeta(`DELETE FROM blogs WHERE id = $1::int`)).
					WithArgs(tc.mockInputArgs...).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(tc.mockError)

				mock.
					ExpectExec(regexp.QuoteMeta(`DELETE FROM comments WHERE blog_id = $1::int`)).
					WithArgs(tc.mockInputArgs...).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(tc.mockError)
			}

			blogService := NewBlogsService(logger, db)

			err = blogService.DeleteBlog(context.TODO(), tc.input)
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
