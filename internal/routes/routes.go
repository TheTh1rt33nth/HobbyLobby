package routes

import (
	"github.com/TheTh1rt33nth/HobbyLobby/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", app.HealthCheck)
	router.Get("/hobby-projects/{projectId}", app.HobbyProjectHandler.GetHobbyProjectById)

	return router
}
