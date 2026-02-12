package main

import (
	"DB/internal/controllers"
	"DB/internal/initializers"
)

func main() {
	// config
	config, err := initializers.LoadConfig("./config.env")
	if err != nil {
		return
	}

	// postgres
	err = initializers.ConnectDB(config)
	if err != nil {
		panic(err)
	}

	// gRPC
	err = initializers.ConnectGRPC(config, &controllers.Server{})
	if err != nil {
		panic(err)
	}
}
