package service_test

import (
	"context"
	"testing"

	v1 "backend/api/v1"
	"backend/internal/model"
	"backend/internal/service"
	mock_repository "backend/test/mocks/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoleRepo := mock_repository.NewMockRoleRepository(ctrl)
	mockMenuRepo := mock_repository.NewMockMenuRepository(ctrl)
	mockAvatarStorage := mock_repository.NewMockAvatarStorage(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)

	userService := service.NewUserService(srv, mockUserRepo, mockRoleRepo, mockMenuRepo, mockAvatarStorage)

	ctx := context.Background()
	req := &v1.RegisterRequest{
		Email:    "test@example.com",
		Password: "123456",
	}

	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(model.User{}, gorm.ErrRecordNotFound)
	mockTm.EXPECT().Transaction(ctx, gomock.Any()).Return(nil)

	err := userService.Register(ctx, req)

	assert.NoError(t, err)
}

func TestUserService_Register_UserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoleRepo := mock_repository.NewMockRoleRepository(ctrl)
	mockMenuRepo := mock_repository.NewMockMenuRepository(ctrl)
	mockAvatarStorage := mock_repository.NewMockAvatarStorage(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	userService := service.NewUserService(srv, mockUserRepo, mockRoleRepo, mockMenuRepo, mockAvatarStorage)

	ctx := context.Background()
	req := &v1.RegisterRequest{
		Email:    "existing@example.com",
		Password: "123456",
	}

	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(model.User{
		Model: gorm.Model{ID: 1},
	}, nil)

	err := userService.Register(ctx, req)

	assert.Error(t, err)
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoleRepo := mock_repository.NewMockRoleRepository(ctrl)
	mockMenuRepo := mock_repository.NewMockMenuRepository(ctrl)
	mockAvatarStorage := mock_repository.NewMockAvatarStorage(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	userService := service.NewUserService(srv, mockUserRepo, mockRoleRepo, mockMenuRepo, mockAvatarStorage)

	ctx := context.Background()
	req := &v1.LoginRequest{
		Username: "testuser",
		Password: "password",
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal("failed to hash password")
	}

	mockUserRepo.EXPECT().GetByUsernameOrEmail(ctx, req.Username, req.Username).Return(model.User{
		Model:    gorm.Model{ID: 1},
		Username: req.Username,
		Password: string(hashedPassword),
	}, nil)

	token, err := userService.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, token)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoleRepo := mock_repository.NewMockRoleRepository(ctrl)
	mockMenuRepo := mock_repository.NewMockMenuRepository(ctrl)
	mockAvatarStorage := mock_repository.NewMockAvatarStorage(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	userService := service.NewUserService(srv, mockUserRepo, mockRoleRepo, mockMenuRepo, mockAvatarStorage)

	ctx := context.Background()
	req := &v1.LoginRequest{
		Username: "nonexistent",
		Password: "password",
	}

	mockUserRepo.EXPECT().GetByUsernameOrEmail(ctx, req.Username, req.Username).Return(model.User{}, gorm.ErrRecordNotFound)

	_, err := userService.Login(ctx, req)

	assert.Error(t, err)
}

func TestUserService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoleRepo := mock_repository.NewMockRoleRepository(ctrl)
	mockMenuRepo := mock_repository.NewMockMenuRepository(ctrl)
	mockAvatarStorage := mock_repository.NewMockAvatarStorage(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	userService := service.NewUserService(srv, mockUserRepo, mockRoleRepo, mockMenuRepo, mockAvatarStorage)

	ctx := context.Background()
	userId := uint(123)

	mockUserRepo.EXPECT().Get(ctx, userId).Return(model.User{
		Model:    gorm.Model{ID: userId},
		Username: "testuser",
	}, nil)
	mockUserRepo.EXPECT().GetRoles(ctx, userId).Return([]string{}, nil)
	mockAvatarStorage.EXPECT().GetURL(ctx, gomock.Any()).Return("", nil)

	user, err := userService.Get(ctx, userId)

	assert.NoError(t, err)
	assert.Equal(t, userId, user.ID)
}
