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
)

func main() {
	//load config
	config := config.MustLoad()
	//database setup
	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /", student.New())

	//setup server
	server := http.Server{
		Addr:    config.Address,
		Handler: router,
	}

	slog.Info("Server Started", slog.String("address", config.Address))

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

	err := server.Shutdown(cntxt)
	if err != nil {
		slog.Info("failed to stop the server", slog.String("error", err.Error()))
	}

}
