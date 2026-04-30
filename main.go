package main

import (
	"log"
	"net/http"
	"os"
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

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
