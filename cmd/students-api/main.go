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

	config "github.com/NazimRiyadh/student_api_golang/internal/config"
	"github.com/NazimRiyadh/student_api_golang/internal/http/handlers/student"
	"github.com/NazimRiyadh/student_api_golang/internal/storage/sqlite"
)

func main() {
	//load config
	config := config.MustLoad()

	//database setup
	storage, err := sqlite.New(config)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Storgae Initialized", slog.String("env", config.Env), slog.String("version", "1.0.0"))

	//setup router
	router := http.NewServeMux()

	//routes
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteById(storage))
	router.HandleFunc("PUT /api/students/{id}", student.UpdateById(storage))

	//setup server
	server := http.Server{
		Addr:    config.Address,
		Handler: router,
	}

	slog.Info("Server Started", slog.String("address", config.Address))

	//Graceful Shutdown Implementation
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server: ", err)
		}
	}()

	<-done

	slog.Info("Shutting down the server")

	cntxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(cntxt)
	if err != nil {
		slog.Info("failed to stop the server", slog.String("error", err.Error()))
	}

}
