package main

import (
	"DB/internal/controllers"
	"DB/internal/initializers"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	// error group
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	// config
	config, err := initializers.LoadConfig("./config.env")
	if err != nil {
		panic(err)
	}

	// postgres
	err = initializers.ConnectDB(config)
	if err != nil {
		panic(err)
	}

	// gRPC
	server := &controllers.Server{}
	g.Go(func() error {
		return initializers.ConnectGRPC(config, server.RegisterAllServices())
	})

	// handle error
	if err := g.Wait(); err != nil {
		log.Fatal("Program terminated: " + err.Error())
	}
}
