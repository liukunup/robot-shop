package repository

import (
	"backend/pkg/log"
	"context"
	"testing"
	"time"

	"backend/internal/constant"
	"backend/internal/model"
	"backend/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var logger *log.Logger

func setupRepository(t *testing.T) (repository.UserRepository, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm connection: %v", err)
	}

	repo := repository.NewRepository(db, nil, nil, nil, nil, logger)
	userRepo := repository.NewUserRepository(repo)

	return userRepo, mock
}

func TestUserRepository_Create(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	user := &model.User{
		Model: gorm.Model{
			ID:        10000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email:    "test@local.lan",
		Username: "test",
		Password: "123456",
		Avatar:   "https://cravatar.cn/avatar/hash?s=100&d=robohash",
		Nickname: "Test",
		Bio:      "It's a test user",
		Language: "zh-CN",
		Timezone: "Asia/Shanghai",
		Theme:    "light",
		Status:   constant.UserStatusActive,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs(user.ID, user.Email, user.Username, user.Password, user.Avatar,
			user.Nickname, user.Bio, user.Language, user.Timezone, user.Theme, user.Status,
			user.CreatedAt, user.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	user := &model.User{
		Model: gorm.Model{
			ID:        10000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email:    "test@local.lan",
		Username: "test",
		Password: "123456",
		Avatar:   "https://cravatar.cn/avatar/hash?s=100&d=robohash",
		Nickname: "Test",
		Bio:      "It's a test user",
		Language: "zh-CN",
		Timezone: "Asia/Shanghai",
		Theme:    "light",
		Status:   constant.UserStatusActive,
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := userRepo.Update(ctx, user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetById(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	userId := "123"

	rows := sqlmock.NewRows([]string{"id", "user_id", "username", "nickname", "password", "email", "created_at", "updated_at"}).
		AddRow(1, "123", "test", "Test", "password", "test@example.com", time.Now(), time.Now())
	mock.ExpectQuery("SELECT \\* FROM `users`").WillReturnRows(rows)

	user, err := userRepo.GetByID(ctx, userId)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "123", user.UserId)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByUsername(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	email := "test@example.com"

	rows := sqlmock.NewRows([]string{"id", "user_id", "username", "nickname", "password", "email", "created_at", "updated_at"}).
		AddRow(1, "123", "test", "Test", "password", "test@example.com", time.Now(), time.Now())
	mock.ExpectQuery("SELECT \\* FROM `users`").WillReturnRows(rows)

	user, err := userRepo.GetByEmail(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)

	assert.NoError(t, mock.ExpectationsWereMet())
}
