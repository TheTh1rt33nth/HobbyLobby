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
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

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
