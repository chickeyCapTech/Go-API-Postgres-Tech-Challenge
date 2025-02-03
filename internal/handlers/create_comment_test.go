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

func TestHandleCreateComment(t *testing.T) {
	tests := map[string]struct {
		wantStatus int
		wantBody   models.Comment
		input      models.Comment
	}{
		"happy path": {
			wantStatus: 200,
			wantBody: models.Comment{
				BlogID:      1,
				UserID:      1,
				Message:     "Good blog",
				CreatedDate: time.Date(2025, 1, 21, 11, 12, 11, 11, time.UTC),
			},
			input: models.Comment{
				BlogID:  1,
				UserID:  1,
				Message: "Good blog",
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a new request
			reqBody, _ := json.Marshal(tc.input)
			req := httptest.NewRequest("POST", "/comments", bytes.NewBuffer(reqBody))

			// Create a new response recorder
			rec := httptest.NewRecorder()

			// Create a new logger
			logger := slog.Default()

			userCreator := new(mock.CommentCreator)
			userCreator.On("CreateComment", context.Background(), tc.input).Return(tc.wantBody, nil)

			// Call the handler
			handler := HandleCreateComment(logger, userCreator)

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
