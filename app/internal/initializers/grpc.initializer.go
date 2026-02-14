package initializers

import (
	"DB/internal/models"
	"log"
	"net"

	"google.golang.org/grpc"
)

func ConnectGRPC(config models.Config, registerFunc func(server *grpc.Server)) error {
	lis, err := net.Listen(config.GrpcNetwork, ":"+config.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	registerFunc(s)

	log.Printf("Server listening at %v", lis.Addr())

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return err
}
