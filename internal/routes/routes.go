package routes

import (
	"log/slog"
	"net/http"

	_ "github.com/chickey/blog/cmd/api/docs"
	"github.com/chickey/blog/internal/handlers"
	"github.com/chickey/blog/internal/services"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title						Blog Service API
// @version					1.0
// @description				Practice Go API using the Standard Library and Postgres
// @termsOfService				http://swagger.io/terms/
// @contact.name				API Support
// @contact.url				http://www.swagger.io/support
// @contact.email				support@swagger.io
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @host						localhost:8000
// @BasePath					/api
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func AddRoutes(mux *http.ServeMux, logger *slog.Logger, usersService *services.UsersService, baseURL string) {
	// Read a user
	mux.Handle("GET /api/users/{id}", handlers.HandleReadUser(logger, usersService))
	mux.Handle("GET /api/users", handlers.HandleListUsers(logger, usersService))
	mux.Handle("POST /api/users", handlers.HandleCreateUser(logger, usersService))
	mux.Handle("PUT /api/users/{id}", handlers.HandleUpdateUser(logger, usersService))
	mux.Handle("DELETE /api/users/{id}", handlers.HandleDeleteUser(logger, usersService))

	// swagger docs
	mux.Handle(
		"GET /swagger/",
		httpSwagger.Handler(httpSwagger.URL(baseURL+"/swagger/doc.json")),
	)
	logger.Info("Swagger running", slog.String("url", baseURL+"/swagger/index.html"))
}
