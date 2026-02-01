package service

import desc "DB/internal/gen/db_v1"

type Server struct {
	desc.UnimplementedDatabaseServiceServer
}

func NewServer() *Server {
	return &Server{}
}
