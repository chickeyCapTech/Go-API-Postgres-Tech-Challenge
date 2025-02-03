package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/chickey/blog/internal/handlers/mock"
	"github.com/chickey/blog/internal/models"
)

func TestHandleDeleteComment(t *testing.T) {
	tests := map[string]struct {
		wantStatus int
		wantBody   models.Comment
		input      models.Comment
	}{
		"happy path": {
			wantStatus: 200,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a new request
			reqBody, _ := json.Marshal(tc.input)
			req := httptest.NewRequest("DELETE", "/comments", bytes.NewBuffer(reqBody))
			req.SetPathValue("id", "1")

			// Create a new response recorder
			rec := httptest.NewRecorder()

			// Create a new logger
			logger := slog.Default()

			userDeleter := new(mock.CommentDeleter)
			userDeleter.On("DeleteComment", context.Background(), uint(0), uint(0)).Return(nil)

			// Call the handler
			handler := HandleDeleteComment(logger, userDeleter)

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
