package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/TheTh1rt33nth/HobbyLobby/internal/api"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/store"
	"github.com/TheTh1rt33nth/HobbyLobby/migrations"
)

type Application struct {
	Logger              *log.Logger
	HobbyProjectHandler *api.HobbyProjectHandler
	UserHandler         *api.UserHandler
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
	userStore := store.NewPostgresUserStore(pgDb)

	// Handlers
	hobbyProjectHandler := api.NewHobbyProjectHandler(hobbyProjectStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)

	app := &Application{
		Logger:              logger,
		HobbyProjectHandler: hobbyProjectHandler,
		UserHandler:         userHandler,
		DB:                  pgDb,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Healthy")
}
