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
	"github.com/ChandanJnv/students-api/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Database initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()
	router.HandleFunc(http.MethodGet+" /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Students API is running"))
	})
	router.HandleFunc(http.MethodPost+" /api/students", student.New(storage))
	router.HandleFunc(http.MethodGet+" /api/students/{id}", student.GetById(storage))
	router.HandleFunc(http.MethodGet+" /api/students/all", student.GetAllStudents(storage))
	router.HandleFunc(http.MethodDelete+" /api/students/{id}", student.DeleteById(storage))
	router.HandleFunc(http.MethodPut+" /api/students/{id}", student.UpdateById(storage))

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
