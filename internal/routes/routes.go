package routes

import (
	"github.com/TheTh1rt33nth/HobbyLobby/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()

	// Populate user context and require authentication for these
	router.Group(func(r chi.Router) {
		r.Use(app.UserMiddleware.PopulateUserContext)
		r.Use(app.UserMiddleware.RequireAuthenticatedUserContext)

		// Hobby projects
		r.Get("/api/hobby-projects/{projectId}", app.HobbyProjectHandler.GetHobbyProjectById)
		r.Post("/api/hobby-projects", app.HobbyProjectHandler.CreateHobbyProject)
		r.Put("/api/hobby-projects/{projectId}", app.HobbyProjectHandler.UpdateHobbyProject)
		r.Delete("/api/hobby-projects/{projectId}", app.HobbyProjectHandler.DeleteHobbyProject)

	})

	router.Get("/health", app.HealthCheck)

	// Users
	router.Post("/api/users/register", app.UserHandler.RegisterUser)

	// Auth & Tokens
	// TODO: dedicated login endpoint that
	router.Post("/api/tokens/auth", app.TokenHandler.CreateToken)

	return router
}
