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
	DB                  *sql.DB
}

func NewApplication(logger *log.Logger) (*Application, error) {

	logger.Println("Connecting to the DB...")

	pgDb, err := store.Open()
	if err != nil {
		panic(err)
	}

	logger.Println("Connected to the DB")

	logger.Println("Migrating the DB...")

	err = store.MigrateFS(pgDb, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger.Println("DB migrated successfully")

	// Stores
	hobbyProjectStore := store.NewPostgresHobbyProjectStore(pgDb)

	// Handlers
	hobbyProjectHandler := api.NewHobbyProjectHandler(hobbyProjectStore)

	app := &Application{
		Logger:              logger,
		HobbyProjectHandler: hobbyProjectHandler,
		DB:                  pgDb,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Healthy")
}
