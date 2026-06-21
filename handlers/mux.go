package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jp-milhomem/first-crud/database"
	"github.com/jp-milhomem/first-crud/utils"
)

// Handler
func NewHandler() http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(utils.SetJSON)

	database := database.Create()

	//Routes

	// Find user
	r.Get("/api/user/{id}", database.FindById())

	//List all users
	r.Get("/api/users", database.FindAll())

	//Create a user
	r.Post("/api/users", database.Insert())

	r.Delete("/api/users/{id}", database.Delete())

	r.Put("/api/users/{id}", database.Update())

	return r
}
