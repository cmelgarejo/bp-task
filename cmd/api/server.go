package main

import (
	"bp-task/internal/database"
	"bp-task/internal/handlers"
	"bp-task/internal/middleware"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx := context.Background()
	dbURL := os.Getenv("DATABASE_URL")
	// Initialize the database
	db, err := database.NewDB(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create and register the / path with the IPFS handler to handle all the requests
	ipfsHandler := middleware.BasicAuthMiddleware(handlers.IPFSHandler(db))

	// Gets host and port from env
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	host := os.Getenv("API_HOST")

	// Set up HTTP server
	mux := http.NewServeMux()
	mux.Handle("/", ipfsHandler)

	// Create http server with timeouts (prevents gosec G112 - Potential DoS vulnerability via resource exhaustion (slowris)
	// https://github.com/securego/gosec#usage
	server := &http.Server{
		Addr:         host + ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second, // Set the read and write timeout to 10 seconds, fast enough
		WriteTimeout: 10 * time.Second,
	}

	// Graceful server shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	// Start the server
	slog.InfoContext(ctx, "Server listening on "+host)

	// Start server
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic("Error starting server: " + err.Error())
	}
}
