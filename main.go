package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Piyuuussshhh/weather-api/api"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// This func will run in parallel to main() and will listen for OS signals.
	// If it receives an interrupt or termination signal, it will cancel the context.
	go func() {
		// Channel for listening to OS signals.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()

	if err := api.Route(ctx); err != nil {
		log.Fatalf("[ERROR] Failed to start server: %v", err)
	}
}