package service

import (
	"DB/internal/initializers"
	"DB/internal/models"
	"context"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/db_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.User, error) {
	user := &models.User{
		Name:     req.GetName(),
		Password: req.GetPassword(),
	}

	result := initializers.DB.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &desc.User{
		Id:   user.ID,
		Name: user.Name,
	}, nil
}

func (s *Server) GetAllUsers(ctx context.Context, req *desc.GetAllUsersRequest) (*desc.GetAllUsersResponse, error) {
	var dbUsers []models.User

	// results with pagination
	result := initializers.DB.Limit(int(req.Limit)).Offset(int(req.Offset)).Find(&dbUsers)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch users: %v", result.Error)
	}

	protoUsers := make([]*desc.User, 0, len(dbUsers))

	for _, u := range dbUsers {
		protoUsers = append(protoUsers, &desc.User{
			Id:   u.ID,
			Name: u.Name,
		})
	}

	return &desc.GetAllUsersResponse{
		Users: protoUsers,
	}, nil
}

func (s *Server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.User, error) {
	var user models.User
	result := initializers.DB.First(user, req.GetId())

	if result.Error != nil {
		return nil, result.Error
	}

	return &desc.User{
		Id:   user.ID,
		Name: user.Name,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*desc.User, error) {
	newUser := &models.User{
		ID:       req.GetId(),
		Name:     req.GetName(),
		Password: req.GetPassword(),
	}

	result := initializers.DB.Save(newUser)

	if result.Error != nil {
		return nil, result.Error
	}

	return &desc.User{
		Id:   newUser.ID,
		Name: newUser.Name,
	}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*desc.DeleteUserResponse, error) {
	result := initializers.DB.Delete(&models.User{}, req.GetId())

	if result.Error != nil {
		return nil, result.Error
	}

	return &desc.DeleteUserResponse{}, nil
}
