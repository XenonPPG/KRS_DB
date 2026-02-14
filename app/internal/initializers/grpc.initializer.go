package initializers

import (
	"DB/internal/models"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func ConnectGRPC(config models.Config, registerFunc func(server *grpc.Server)) error {
	lis, err := net.Listen(config.GrpcNetwork, ":"+config.GrpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer()
	registerFunc(s)

	// server error channel
	serverError := make(chan error, 1)

	// serve
	go func() {
		log.Printf("Server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			serverError <- err
		}
	}()

	// system signals channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// block until error or signal
	select {
	case err := <-serverError:
		return fmt.Errorf("gRPC server error: %w", err)
	case sig := <-stop:
		log.Printf("Received signal: %v. Shutting down...", sig)
		s.GracefulStop()
		log.Println("gRPC server gracefully stopped")
		return nil
	}
}
