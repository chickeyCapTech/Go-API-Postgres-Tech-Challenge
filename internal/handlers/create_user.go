package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/chickey/blog/internal/models"
)

// userCreator represents a type capable of reading a user from storage and
// returning it or an error.
type userCreator interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
}

// createUserRequest represents the request for creating a user.
type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// createUserResponse represents the response for creating a user.
type createUserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary		Create User
// @Description	Creates a User
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			request	body		createUserRequest	true	"User to Create"
// @Success		200		{object}	uint
// @Failure		400		{object}	string
// @Failure		404		{object}	string
// @Failure		500		{object}	string
// @Router			/users  [POST]
func HandleCreateUser(logger *slog.Logger, userCreator userCreator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var request models.User
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to decode request",
				slog.String("error", err.Error()))

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		// Read the user
		user, err := userCreator.CreateUser(ctx, request)
		if err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to create user",
				slog.String("error", err.Error()),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Convert our models.User domain model into a response model.
		response := createUserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
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
