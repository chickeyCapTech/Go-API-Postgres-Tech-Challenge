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
func AddRoutes(mux *http.ServeMux, logger *slog.Logger, usersService *services.UsersService, blogsService *services.BlogsService, commentsService *services.CommentsService, baseURL string) {
	// User endpoints
	mux.Handle("GET /api/user/{id}", handlers.HandleReadUser(logger, usersService))
	mux.Handle("GET /api/user", handlers.HandleListUsers(logger, usersService))
	mux.Handle("POST /api/user", handlers.HandleCreateUser(logger, usersService))
	mux.Handle("PUT /api/user/{id}", handlers.HandleUpdateUser(logger, usersService))
	mux.Handle("DELETE /api/user/{id}", handlers.HandleDeleteUser(logger, usersService))

	// Blog endpoints
	mux.Handle("GET /api/blog/{id}", handlers.HandleReadBlog(logger, blogsService))
	mux.Handle("GET /api/blog", handlers.HandleListBlogs(logger, blogsService))
	mux.Handle("POST /api/blog", handlers.HandleCreateBlog(logger, blogsService))
	mux.Handle("PUT /api/blog/{id}", handlers.HandleUpdateBlog(logger, blogsService))
	mux.Handle("DELETE /api/blog/{id}", handlers.HandleDeleteBlog(logger, blogsService))

	// Comment endpoints
	mux.Handle("GET /api/comment", handlers.HandleListComments(logger, commentsService))
	mux.Handle("POST /api/comment", handlers.HandleCreateComment(logger, commentsService))
	mux.Handle("PUT /api/comment", handlers.HandleUpdateComment(logger, commentsService))
	mux.Handle("DELETE /api/comment", handlers.HandleDeleteComment(logger, commentsService))

	// health check
	mux.Handle("GET /api/health", handlers.HandleHealthCheck(logger))

	// swagger docs
	mux.Handle(
		"GET /swagger/",
		httpSwagger.Handler(httpSwagger.URL(baseURL+"/swagger/doc.json")),
	)
	logger.Info("Swagger running", slog.String("url", baseURL+"/swagger/index.html"))
}
