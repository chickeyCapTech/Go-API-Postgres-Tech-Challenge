package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/chickey/blog/internal/models"
)

// userReader represents a type capable of reading a user from storage and
// returning it or an error.
type usersReader interface {
	ListUsers(ctx context.Context) ([]models.User, error)
}

// listUsersResponse represents the response for listing users.
type listUsersResponse struct {
	Users []readUserResponse
}

// @Summary		List Users
// @Description	List All Users
// @Tags			users
// @Accept			json
// @Produce		json
// @Success		200	{array}		models.User
// @Failure		400	{object}	string
// @Failure		404	{object}	string
// @Failure		500	{object}	string
// @Router			/users  [GET]
func HandleListUsers(logger *slog.Logger, usersReader usersReader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Read the user
		users, err := usersReader.ListUsers(ctx)
		if err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to list users",
				slog.String("error", err.Error()),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Convert our models.User domain model into a response model.
		response := listUsersResponse{
			Users: []readUserResponse{},
		}

		for _, user := range users {
			newUser := readUserResponse{
				ID:       user.ID,
				Name:     user.Name,
				Email:    user.Email,
				Password: user.Password,
			}
			response.Users = append(response.Users, newUser)
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
