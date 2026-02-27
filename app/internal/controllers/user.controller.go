package controllers

import (
	"DB/internal/initializers"
	"context"
	"errors"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/user_v1"
	"github.com/XenonPPG/KRS_CONTRACTS/models"
	"golang.org/x/crypto/bcrypt"
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
		Role:       req.GetRole(),
		ColorTheme: req.GetColorTheme(),
	}

	result := initializers.DB.Create(user)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", result.Error)
	}

	return &desc.User{
		Id:         user.ID,
		Login:      user.Login,
		Role:       &user.Role,
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

	protoUsers := make([]*desc.User, 0, len(dbUsers))

	for _, u := range dbUsers {
		protoUsers = append(protoUsers, &desc.User{
			Id:         u.ID,
			Login:      u.Login,
			Role:       &u.Role,
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
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	colorTheme := user.ColorTheme
	return &desc.User{
		Id:         user.ID,
		Login:      user.Login,
		Role:       &user.Role,
		ColorTheme: &colorTheme,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*desc.User, error) {
	// get old user
	oldUser, err := s.GetUser(ctx, &desc.GetUserRequest{
		Id: req.GetId(),
	})
	if err != nil || oldUser == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	updatedUser := &models.User{
		ID:         req.GetId(),
		Login:      req.GetLogin(),
		Role:       req.GetRole(),
		ColorTheme: req.GetColorTheme(),
	}

	// check if something was updated
	anythingUpdated := oldUser.Login != updatedUser.Login
	anythingUpdated = anythingUpdated || oldUser.ColorTheme != &updatedUser.ColorTheme
	anythingUpdated = anythingUpdated || oldUser.Role != &updatedUser.Role
	if !anythingUpdated {
		return nil, status.Errorf(codes.InvalidArgument, "nothing to update")
	}

	result := initializers.DB.Model(&models.User{ID: req.GetId()}).Updates(updatedUser)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", result.Error)
	}

	return &desc.User{
		Id:         updatedUser.ID,
		Login:      updatedUser.Login,
		Role:       &updatedUser.Role,
		ColorTheme: &updatedUser.ColorTheme,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *desc.LoginRequest) (*desc.User, error) {
	// try to find the user
	var user models.User
	result := initializers.DB.First(&user, "login = ?", req.GetLogin())
	if result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	// check password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.GetPassword()))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password")
	}

	// return user
	return &desc.User{
		Id:         user.ID,
		Login:      user.Login,
		Role:       &user.Role,
		ColorTheme: &user.ColorTheme,
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
	// get user
	user, err := s.GetUser(ctx, &desc.GetUserRequest{Id: req.GetId()})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	// check if not admin
	if user.GetRole() == desc.UserRole_ADMIN {
		return nil, status.Errorf(codes.PermissionDenied, "cannot delete admin")
	}

	// delete
	result := initializers.DB.Delete(&models.User{}, req.GetId())

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.GetId())
	}

	return &emptypb.Empty{}, nil
}
