package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

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

	const maxRetries = 5
	for i := 0; i < maxRetries; i++ {
		if err = pgDb.PingContext(context.Background()); err == nil {
			break
		}
		if i == maxRetries-1 {
			pgDb.Close()
			return nil, fmt.Errorf("database not reachable after %d attempts: %w", maxRetries, err)
		}
		wait := time.Duration(1<<uint(i)) * time.Second
		logger.Printf("DB not ready (attempt %d/%d): %v. Retrying in %s...", i+1, maxRetries, err, wait)
		time.Sleep(wait)
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
	if err := app.DB.PingContext(r.Context()); err != nil {
		app.Logger.Printf("HealthCheck: DB ping failed: %v", err)
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}
	fmt.Fprint(w, "Healthy")
}
