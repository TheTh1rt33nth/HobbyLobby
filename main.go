package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TheTh1rt33nth/HobbyLobby/internal/app"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/routes"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/store"
	"github.com/TheTh1rt33nth/HobbyLobby/migrations"
	"github.com/cenkalti/backoff/v5"
)

func runMigrations(logger *log.Logger) error {
	logger.Println("Connecting to DB for migrations...")

	db, err := store.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	be := backoff.NewExponentialBackOff()
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 30 * time.Second

	if _, err = backoff.Retry(
		context.Background(),
		func() (struct{}, error) { return struct{}{}, db.PingContext(context.Background()) },
		backoff.WithBackOff(be),
		backoff.WithNotify(func(err error, d time.Duration) {
			logger.Printf("DB not ready: %v. Retrying in %s...", err, d)
		}),
	); err != nil {
		return err
	}

	logger.Println("Running migrations...")
	if err := store.MigrateFS(db, migrations.FS, "."); err != nil {
		return err
	}

	logger.Println("Migrations completed successfully")
	return nil
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		if err := runMigrations(logger); err != nil {
			logger.Fatalf("Migration failed: %v", err)
		}
		return
	}

	app, err := app.NewApplication(logger)
	if err != nil {
		panic(err)
	}

	defer app.DB.Close()

	app.Logger.Println("We're alive")

	router := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Give in-flight requests 30s to finish before the process exits
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown with dignity
	defer func() {
		if err := server.Shutdown(ctx); err != nil {
			app.Logger.Printf("Server forced to shutdown: %v", err)
		}

		app.Logger.Println("Server stopped")
	}()

	// Buffer of 1 so the signal sender is never blocked even if we haven't called <-quit yet.
	quit := make(chan os.Signal, 1)
	errs := make(chan error, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, os.Interrupt)

	go func() {
		errs <- server.ListenAndServe()
	}()

	for {
		select {
		case <-quit:
			signal.Stop(quit)
			app.Logger.Println("Shutdown signal received, draining connections...")
			return
		case err = <-errs:
			if err != nil {
				app.Logger.Fatal(err)
			}
			signal.Stop(quit)
		}
	}
}
