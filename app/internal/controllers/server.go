package controllers

import (
	noteDesc "github.com/XenonPPG/KRS_CONTRACTS/gen/note_v1"
	userDesc "github.com/XenonPPG/KRS_CONTRACTS/gen/user_v1"
	"google.golang.org/grpc"
)

type Server struct {
	userDesc.UnimplementedUserServiceServer
	noteDesc.UnimplementedNoteServiceServer
}

func (s *Server) RegisterAllServices() func(*grpc.Server) {
	return func(server *grpc.Server) {
		userDesc.RegisterUserServiceServer(server, s)
		noteDesc.RegisterNoteServiceServer(server, s)
	}
}
