package initializers

import (
	"DB/internal/models"
	"log"
	"net"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/db_v1"
	"google.golang.org/grpc"
)

func ConnectGRPC(config models.Config, server desc.DatabaseServiceServer) error {
	lis, err := net.Listen(config.GrpcNetwork, ":"+config.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	desc.RegisterDatabaseServiceServer(s, server)

	log.Printf("Server listening at %v", lis.Addr())

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return err
}
