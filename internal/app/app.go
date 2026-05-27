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
	"github.com/cenkalti/backoff/v5"
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

	be := backoff.NewExponentialBackOff()
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 30 * time.Second

	if _, err = backoff.Retry(
		context.Background(),
		func() (struct{}, error) { return struct{}{}, pgDb.PingContext(context.Background()) },
		backoff.WithBackOff(be),
		backoff.WithNotify(func(err error, d time.Duration) {
			logger.Printf("DB not ready: %v. Retrying in %s...", err, d)
		}),
	); err != nil {
		pgDb.Close()
		return nil, fmt.Errorf("database not reachable: %w", err)
	}

	logger.Println("Connected to the DB")

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

	app.startTokenCleanup(time.Hour)

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

// ReadinessCheck is the Kubernetes readiness probe (/ready).
// Returns 200 only when the database is reachable
func (app *Application) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	if err := app.DB.PingContext(r.Context()); err != nil {
		app.Logger.Printf("ReadinessCheck: DB ping failed: %v", err)
		http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		return
	}
	fmt.Fprint(w, "Ready")
}

// startTokenCleanup runs a background goroutine that calls the
// cleanup_expired_tokens() Postgres function on the given interval
func (app *Application) startTokenCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if _, err := app.DB.Exec(`SELECT cleanup_expired_tokens()`); err != nil {
				app.Logger.Printf("token cleanup failed: %v", err)
			}
		}
	}()
}
