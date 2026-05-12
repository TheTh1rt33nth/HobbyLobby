package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/TheTh1rt33nth/HobbyLobby/internal/api"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/middleware"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/store"
	"github.com/TheTh1rt33nth/HobbyLobby/migrations"
)

type Application struct {
	Logger              *log.Logger
	HobbyProjectHandler *api.HobbyProjectHandler
	PaintProfileHandler *api.PaintProfileHandler
	UserHandler         *api.UserHandler
	TokenHandler        *api.TokenHandler
	UserMiddleware      middleware.UserMiddleware
	DB                  *sql.DB
}

func NewApplication(logger *log.Logger) (*Application, error) {

	logger.Println("Connecting to the DB...")

	pgDb, err := store.Open()
	if err != nil {
		return nil, err
	}

	logger.Println("Connected to the DB")

	logger.Println("Migrating the DB...")

	err = store.MigrateFS(pgDb, migrations.FS, ".")
	if err != nil {
		return nil, err
	}

	logger.Println("DB migrated successfully")

	// Stores
	hobbyProjectStore := store.NewPostgresHobbyProjectStore(pgDb)
	paintProfileStore := store.NewPostgresPaintProfileStore(pgDb)
	userStore := store.NewPostgresUserStore(pgDb)
	tokenStore := store.NewPostgresTokenStore(pgDb)

	// Handlers
	hobbyProjectHandler := api.NewHobbyProjectHandler(hobbyProjectStore, paintProfileStore, logger)
	paintProfileHandler := api.NewPaintProfileHandler(paintProfileStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

	// Middleware
	userMiddleware := middleware.UserMiddleware{UserStore: userStore}

	app := &Application{
		Logger:              logger,
		HobbyProjectHandler: hobbyProjectHandler,
		PaintProfileHandler: paintProfileHandler,
		UserHandler:         userHandler,
		TokenHandler:        tokenHandler,
		UserMiddleware:      userMiddleware,
		DB:                  pgDb,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Healthy")
}
