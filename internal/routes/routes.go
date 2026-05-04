package routes

import (
	"github.com/TheTh1rt33nth/HobbyLobby/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", app.HealthCheck)

	// Hobby projects
	router.Get("/api/hobby-projects/{projectId}", app.HobbyProjectHandler.GetHobbyProjectById)
	router.Post("/api/hobby-projects", app.HobbyProjectHandler.CreateHobbyProject)
	router.Put("/api/hobby-projects/{projectId}", app.HobbyProjectHandler.UpdateHobbyProject)
	router.Delete("/api/hobby-projects/{projectId}", app.HobbyProjectHandler.DeleteHobbyProject)

	return router
}
