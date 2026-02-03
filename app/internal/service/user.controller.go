package service

import (
	"DB/internal/initializers"
	"DB/internal/models"
	"context"
	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/db_v1"
)

func (s *Server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.UserResponse, error) {
	user := &models.User{
		Name:     req.GetName(),
		Password: req.GetPassword(),
	}

	result := initializers.DB.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &desc.UserResponse{
		Id:   user.ID,
		Name: user.Name,
	}, nil
}

func (s *Server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.UserResponse, error) {
	var user models.User
	result := initializers.DB.First(user, req.GetId())

	if result.Error != nil {
		return nil, result.Error
	}

	return &desc.UserResponse{
		Id:   user.ID,
		Name: user.Name,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*desc.UserResponse, error) {
	newUser := &models.User{
		ID:       req.GetId(),
		Name:     req.GetName(),
		Password: req.GetPassword(),
	}

	result := initializers.DB.Save(newUser)

	if result.Error != nil {
		return nil, result.Error
	}

	return &desc.UserResponse{
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
