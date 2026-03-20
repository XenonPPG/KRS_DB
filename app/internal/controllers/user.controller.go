package controllers

import (
	"DB/internal/initializers"
	"context"
	"errors"
	"fmt"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/user_v1"
	"github.com/XenonPPG/KRS_CONTRACTS/models"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
	}, nil
}

func (s *Server) GetAllUsers(ctx context.Context, req *desc.GetAllUsersRequest) (*desc.GetAllUsersResponse, error) {
	var dbUsers []models.User

	order := "DESC"
	if req.GetAscendingOrder() {
		order = "ASC"
	}

	// results with pagination
	result := initializers.DB.
		Order(fmt.Sprintf("id %s", order)).
		Limit(int(req.GetLimit())).
		Offset(int(req.GetOffset())).
		Find(&dbUsers)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch users: %v", result.Error)
	}

	// total user count
	var usersAmount int64
	result = initializers.DB.Model(&models.User{}).Count(&usersAmount)
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
			CreatedAt:  timestamppb.New(u.CreatedAt),
			UpdatedAt:  timestamppb.New(u.UpdatedAt),
		})
	}

	return &desc.GetAllUsersResponse{
		Users:      protoUsers,
		TotalCount: int32(usersAmount),
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
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
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
		Login:      req.GetLogin(),
		Role:       req.GetRole(),
		ColorTheme: req.GetColorTheme(),
	}

	// check if something was updated
	updateMap := make(map[string]any)
	if oldUser.GetLogin() != updatedUser.Login {
		updateMap["login"] = updatedUser.Login
	}
	if oldUser.GetRole() != updatedUser.Role {
		updateMap["role"] = updatedUser.Role
	}
	if oldUser.GetColorTheme() != updatedUser.ColorTheme {
		updateMap["color_theme"] = updatedUser.ColorTheme
	}
	if len(updateMap) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nothing to update")
	}

	result := initializers.DB.Model(&models.User{}).
		Where("id = ?", req.GetId()).
		Updates(updateMap)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", result.Error)
	}

	return &desc.User{
		Id:         updatedUser.ID,
		Login:      updatedUser.Login,
		Role:       &updatedUser.Role,
		ColorTheme: &updatedUser.ColorTheme,
		CreatedAt:  timestamppb.New(updatedUser.CreatedAt),
		UpdatedAt:  timestamppb.New(updatedUser.UpdatedAt),
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
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
	}, nil
}

func (s *Server) UpdatePassword(ctx context.Context, req *desc.UpdatePasswordRequest) (*emptypb.Empty, error) {
	// check if the new password is different
	if req.GetNewPassword() == req.GetOldPassword() {
		return nil, status.Error(codes.InvalidArgument, "new password must be different")
	}

	// get user
	var user models.User
	if err := initializers.DB.First(&user, req.GetId()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, err
	}

	// check if the old password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.GetOldPassword())); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid old password")
	}

	// update
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.GetNewPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	if err = initializers.DB.Model(&user).Update("password", string(bytes)).Error; err != nil {
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
