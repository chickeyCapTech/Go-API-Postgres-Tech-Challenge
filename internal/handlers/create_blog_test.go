package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chickey/blog/internal/handlers/mock"
	"github.com/chickey/blog/internal/models"
)

func TestHandleCreateBlog(t *testing.T) {
	tests := map[string]struct {
		wantStatus int
		wantBody   models.Blog
		input      models.Blog
	}{
		"happy path": {
			wantStatus: 200,
			wantBody: models.Blog{
				ID:          1,
				AuthorID:    1,
				Title:       "Book Title",
				Score:       8.2,
				CreatedDate: time.Date(2025, 1, 21, 11, 12, 11, 11, time.UTC),
			},
			input: models.Blog{
				AuthorID: 1,
				Title:    "Book Title",
				Score:    8.2,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a new request
			reqBody, _ := json.Marshal(tc.input)
			req := httptest.NewRequest("POST", "/blogs", bytes.NewBuffer(reqBody))

			// Create a new response recorder
			rec := httptest.NewRecorder()

			// Create a new logger
			logger := slog.Default()

			userCreator := new(mock.BlogCreator)
			userCreator.On("CreateBlog", context.Background(), tc.input).Return(tc.wantBody, nil)

			// Call the handler
			handler := HandleCreateBlog(logger, userCreator)

			handler.ServeHTTP(rec, req)
			// Check the status code
			if rec.Code != tc.wantStatus {
				t.Errorf("want status %d, got %d", tc.wantStatus, rec.Code)
			}

			// // Check the body
			// if strings.Trim(rec.Body.String(), "\n") != fmt.Sprintf("%+v", tc.wantBody) {
			// 	t.Errorf("want body %q, got %q", tc.wantBody, rec.Body.String())
			// }
		})
	}
}
