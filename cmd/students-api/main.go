package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChandanJnv/students-api/internal/config"
	"github.com/ChandanJnv/students-api/internal/handlers/student"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// database setup
	// setup router
	router := http.NewServeMux()
	router.HandleFunc(http.MethodPost+" /api/students", student.New())

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("server starting at:", slog.String("address", server.Addr))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("failed to start server: ", err)
		}
	}()

	<-done
	slog.Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server: ", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successful")

}
