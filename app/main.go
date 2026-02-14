package main

import (
	"DB/internal/controllers"
	"DB/internal/initializers"
)

func main() {
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
	err = initializers.ConnectGRPC(config, server.RegisterAllServices())
	if err != nil {
		panic(err)
	}

	select {}
}
