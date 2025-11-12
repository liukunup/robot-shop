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
	"gorm.io/gorm"
)

func TestApiService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApiRepo := mock_repository.NewMockApiRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	apiService := service.NewApiService(srv, mockApiRepo)

	ctx := context.Background()
	req := &v1.ApiSearchRequest{
		Page:     1,
		PageSize: 10,
	}

	mockApiRepo.EXPECT().List(ctx, req).Return([]model.Api{
		{Model: gorm.Model{ID: 1}, Group: "User", Name: "ListUsers"},
		{Model: gorm.Model{ID: 2}, Group: "Robot", Name: "ListRobots"},
	}, int64(2), nil)

	result, err := apiService.List(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), result.Total)
	assert.Len(t, result.List, 2)
}

func TestApiService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApiRepo := mock_repository.NewMockApiRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	apiService := service.NewApiService(srv, mockApiRepo)

	ctx := context.Background()
	req := &v1.ApiRequest{
		Group:  "Test",
		Name:   "TestApi",
		Path:   "/v1/test",
		Method: "POST",
	}

	mockApiRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	err := apiService.Create(ctx, req)

	assert.NoError(t, err)
}

func TestApiService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApiRepo := mock_repository.NewMockApiRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	apiService := service.NewApiService(srv, mockApiRepo)

	ctx := context.Background()
	apiId := uint(1)
	req := &v1.ApiRequest{
		Group:  "Test",
		Name:   "UpdatedApi",
		Path:   "/v1/updated",
		Method: "PUT",
	}

	mockApiRepo.EXPECT().Update(ctx, apiId, gomock.Any()).Return(nil)

	err := apiService.Update(ctx, apiId, req)

	assert.NoError(t, err)
}

func TestApiService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApiRepo := mock_repository.NewMockApiRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	apiService := service.NewApiService(srv, mockApiRepo)

	ctx := context.Background()
	apiId := uint(1)

	mockApiRepo.EXPECT().Delete(ctx, apiId).Return(nil)

	err := apiService.Delete(ctx, apiId)

	assert.NoError(t, err)
}

// 错误处理测试
func TestApiService_Create_DuplicateApi(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApiRepo := mock_repository.NewMockApiRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	apiService := service.NewApiService(srv, mockApiRepo)

	ctx := context.Background()
	req := &v1.ApiRequest{
		Group:  "Test",
		Name:   "TestApi",
		Path:   "/v1/test",
		Method: "POST",
	}

	mockApiRepo.EXPECT().Create(ctx, gomock.Any()).Return(gorm.ErrDuplicatedKey)

	err := apiService.Create(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrDuplicatedKey, err)
}

func TestApiService_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApiRepo := mock_repository.NewMockApiRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	apiService := service.NewApiService(srv, mockApiRepo)

	ctx := context.Background()
	apiId := uint(999)
	req := &v1.ApiRequest{
		Group:  "Test",
		Name:   "UpdatedApi",
		Path:   "/v1/updated",
		Method: "PUT",
	}

	mockApiRepo.EXPECT().Update(ctx, apiId, gomock.Any()).Return(gorm.ErrRecordNotFound)

	err := apiService.Update(ctx, apiId, req)

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestApiService_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApiRepo := mock_repository.NewMockApiRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	apiService := service.NewApiService(srv, mockApiRepo)

	ctx := context.Background()
	apiId := uint(999)

	mockApiRepo.EXPECT().Delete(ctx, apiId).Return(gorm.ErrRecordNotFound)

	err := apiService.Delete(ctx, apiId)

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
