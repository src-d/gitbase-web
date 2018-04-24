package server

import (
	"database/sql"
	"net/http"

	"github.com/src-d/gitbase-playground/server/handler"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pressly/lg"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

// Router returns a Handler to serve the backend
func Router(
	logger *logrus.Logger,
	static *handler.Static,
	version string,
	db *sql.DB,
) http.Handler {

	// cors options
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Location", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(cors.New(corsOptions).Handler)
	r.Use(lg.RequestLogger(logger))

	r.Post("/query", handler.APIHandlerFunc(handler.Query(db)))

	r.Get("/version", handler.APIHandlerFunc(handler.Version(version)))

	r.Get("/static/*", static.ServeHTTP)
	r.Get("/*", static.ServeHTTP)

	return r
}
