package service_test

import (
	"context"
	"testing"

	v1 "backend/api/v1"
	"backend/internal/model"
	"backend/internal/service"
	"backend/test/mocks/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRobotService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRobotRepo := mock_repository.NewMockRobotRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	robotService := service.NewRobotService(srv, mockRobotRepo)

	ctx := context.Background()
	req := &v1.RobotSearchRequest{
		Page:     1,
		PageSize: 10,
	}

	mockRobots := []model.Robot{
		{
			Model:   gorm.Model{ID: 1},
			Name:    "Robot 1",
			Owner:   "user1",
			Enabled: true,
		},
		{
			Model:   gorm.Model{ID: 2},
			Name:    "Robot 2",
			Owner:   "user2",
			Enabled: true,
		},
	}

	mockRobotRepo.EXPECT().List(ctx, req).Return(mockRobots, int64(2), nil)

	result, err := robotService.List(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.Total)
	assert.Len(t, result.List, 2)
	assert.Equal(t, "Robot 1", result.List[0].Name)
}

func TestRobotService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRobotRepo := mock_repository.NewMockRobotRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	robotService := service.NewRobotService(srv, mockRobotRepo)

	ctx := context.Background()
	robotId := uint(1)

	mockRobot := model.Robot{
		Model:   gorm.Model{ID: 1},
		Name:    "Test Robot",
		Owner:   "testuser",
		Enabled: true,
	}

	mockRobotRepo.EXPECT().Get(ctx, robotId).Return(mockRobot, nil)

	result, err := robotService.Get(ctx, robotId)

	assert.NoError(t, err)
	assert.Equal(t, "Test Robot", result.Name)
	assert.Equal(t, "testuser", result.Owner)
}

func TestRobotService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRobotRepo := mock_repository.NewMockRobotRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	robotService := service.NewRobotService(srv, mockRobotRepo)

	ctx := context.Background()
	req := &v1.RobotRequest{
		Name:    "New Robot",
		Desc:    "Test Description",
		Webhook: "https://example.com/webhook",
		Enabled: true,
		Owner:   "testuser",
	}

	mockRobotRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	err := robotService.Create(ctx, req)

	assert.NoError(t, err)
}

func TestRobotService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRobotRepo := mock_repository.NewMockRobotRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	robotService := service.NewRobotService(srv, mockRobotRepo)

	ctx := context.Background()
	robotId := uint(1)
	req := &v1.RobotRequest{
		Name:    "Updated Robot",
		Desc:    "Updated Description",
		Webhook: "https://example.com/webhook-updated",
		Enabled: true,
		Owner:   "updateduser",
	}

	mockRobotRepo.EXPECT().Update(ctx, robotId, gomock.Any()).Return(nil)

	err := robotService.Update(ctx, robotId, req)

	assert.NoError(t, err)
}

func TestRobotService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRobotRepo := mock_repository.NewMockRobotRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	robotService := service.NewRobotService(srv, mockRobotRepo)

	ctx := context.Background()
	robotId := uint(1)

	mockRobotRepo.EXPECT().Delete(ctx, robotId).Return(nil)

	err := robotService.Delete(ctx, robotId)

	assert.NoError(t, err)
}
