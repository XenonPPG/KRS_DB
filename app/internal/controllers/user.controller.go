package controllers

import (
	"DB/internal/initializers"
	"context"
	"errors"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/user_v1"
	"github.com/XenonPPG/KRS_CONTRACTS/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (s *Server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.User, error) {
	_, ok := desc.ColorTheme_name[int32(req.GetColorTheme())]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid color theme")
	}

	user := &models.User{
		Login:      req.GetLogin(),
		Password:   req.GetPassword(),
		ColorTheme: req.GetColorTheme(),
	}

	result := initializers.DB.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &desc.User{
		Id:         user.ID,
		Login:      user.Login,
		ColorTheme: &user.ColorTheme,
	}, nil
}

func (s *Server) GetAllUsers(ctx context.Context, req *desc.GetAllUsersRequest) (*desc.GetAllUsersResponse, error) {
	var dbUsers []models.User

	// results with pagination
	result := initializers.DB.Limit(int(req.Limit)).Offset(int(req.Offset)).Find(&dbUsers)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch users: %v", result.Error)
	}

	protoUsers := make([]*desc.User, len(dbUsers))

	for _, u := range dbUsers {
		protoUsers = append(protoUsers, &desc.User{
			Id:         u.ID,
			Login:      u.Login,
			ColorTheme: &u.ColorTheme,
		})
	}

	return &desc.GetAllUsersResponse{
		Users: protoUsers,
	}, nil
}

func (s *Server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.User, error) {
	var user models.User
	result := initializers.DB.First(&user, req.GetId())

	if result.Error != nil {
		return nil, result.Error
	}

	colorTheme := desc.ColorTheme(user.ColorTheme)
	return &desc.User{
		Id:         user.ID,
		Login:      user.Login,
		ColorTheme: &colorTheme,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*desc.User, error) {
	// get old user
	oldUser, err := s.GetUser(ctx, &desc.GetUserRequest{
		Id: req.GetId(),
	})
	if err != nil {
		return nil, err
	}

	updatedUser := &models.User{
		ID:         req.GetId(),
		Login:      req.GetLogin(),
		ColorTheme: req.GetColorTheme(),
	}

	// check if the login is updated
	if updatedUser.Login == oldUser.GetLogin() {
		return nil, status.Errorf(codes.InvalidArgument, "new login is the same as old one")
	}

	result := initializers.DB.Model(oldUser).Updates(updatedUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &desc.User{
		Id:         updatedUser.ID,
		Login:      updatedUser.Login,
		ColorTheme: &updatedUser.ColorTheme,
	}, nil
}

func (s *Server) VerifyPassword(ctx context.Context, req *desc.VerifyPasswordRequest) (*desc.IsValidResponse, error) {
	// get user
	var user models.User
	result := initializers.DB.First(&user, req.GetId())

	if result.Error != nil {
		return nil, result.Error
	}

	// check hash
	isValid := user.Password == req.GetPassword()
	return &desc.IsValidResponse{
		Valid: isValid,
	}, nil
}

func (s *Server) UpdatePassword(ctx context.Context, req *desc.UpdatePasswordRequest) (*emptypb.Empty, error) {
	// get user
	var user models.User
	if err := initializers.DB.First(&user, req.GetId()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, err
	}

	// check hash
	if user.Password != req.GetOldPassword() {
		return nil, status.Error(codes.InvalidArgument, "incorrect old password")
	}

	// check if the new password is different
	if user.Password == req.GetNewPassword() {
		return nil, status.Error(codes.InvalidArgument, "new password must be different")
	}

	// update
	if err := initializers.DB.Model(&user).Update("password", req.GetNewPassword()).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update password: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	result := initializers.DB.Delete(&models.User{}, req.GetId())

	if result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}
