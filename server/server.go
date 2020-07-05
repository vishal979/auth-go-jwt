package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/vishal979/auth/cmd/filehandler"
	"github.com/vishal979/auth/server/controller"
)

// Run single point of contact for main.go file
func Run() {
	var server = &controller.Server{}
	filehandler.Open()
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Error("Error Loading environment file ", err)
	} else {
		log.Info("environment variables loaded successfully")
	}
	(*server).Initialize()
	srv := &http.Server{
		Handler:      server.Router,
		Addr:         ":7070",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
