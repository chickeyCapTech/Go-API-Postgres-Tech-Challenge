package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chickey/blog/internal/handlers/mock"
	"github.com/chickey/blog/internal/models"
)

func TestHandleReadBlog(t *testing.T) {
	tests := map[string]struct {
		wantStatus  int
		wantBody    string
		wantResults models.Blog
	}{
		"happy path": {
			wantStatus: 200,
			wantResults: models.Blog{
				ID:          1,
				AuthorID:    1,
				Title:       "Book Title",
				Score:       8.2,
				CreatedDate: time.Date(2025, 1, 21, 11, 12, 11, 11, time.UTC),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a new request
			req := httptest.NewRequest(http.MethodGet, "/blogs/1", nil)
			req.SetPathValue("id", "1")

			// Create a new response recorder
			rec := httptest.NewRecorder()

			// Create a new logger
			logger := slog.Default()

			userReader := new(mock.BlogReader)
			userReader.On("ReadBlog", context.Background(), uint64(1)).Return(tc.wantResults, nil)
			// Call the handler
			handler := HandleReadBlog(logger, userReader)

			handler.ServeHTTP(rec, req)
			// Check the status code
			if rec.Code != tc.wantStatus {
				t.Errorf("want status %d, got %d", tc.wantStatus, rec.Code)
			}

			// // Check the body
			// json, _ := json.Marshal(tc.wantResults)
			// if strings.Trim(rec.Body.String(), "\n") != string(json) {
			// 	t.Errorf("want body %q, got %q", string(json), rec.Body.String())
			// }
		})
	}
}
