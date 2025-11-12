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

func TestMenuService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMenuRepo := mock_repository.NewMockMenuRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	menuService := service.NewMenuService(srv, mockMenuRepo)

	ctx := context.Background()
	req := &v1.MenuSearchRequest{
		Page:     1,
		PageSize: 10,
	}

	mockMenus := []model.Menu{
		{
			Model:     gorm.Model{ID: 1},
			Name:      "Dashboard",
			Path:      "/dashboard",
			Component: "@/pages/Dashboard",
			Icon:      "dashboard",
		},
		{
			Model:     gorm.Model{ID: 2},
			Name:      "Admin",
			Path:      "/admin",
			Component: "@/pages/Admin",
			Icon:      "crown",
		},
	}

	mockMenuRepo.EXPECT().List(ctx, req).Return(mockMenus, int64(2), nil)

	result, err := menuService.List(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.Total)
	assert.Len(t, result.List, 2)
	assert.Equal(t, "Dashboard", result.List[0].Name)
}

func TestMenuService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMenuRepo := mock_repository.NewMockMenuRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	menuService := service.NewMenuService(srv, mockMenuRepo)

	ctx := context.Background()
	req := &v1.MenuRequest{
		Name:      "New Menu",
		Path:      "/new",
		Component: "@/pages/New",
		Icon:      "star",
	}

	mockMenuRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	err := menuService.Create(ctx, req)

	assert.NoError(t, err)
}

func TestMenuService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMenuRepo := mock_repository.NewMockMenuRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	menuService := service.NewMenuService(srv, mockMenuRepo)

	ctx := context.Background()
	menuId := uint(1)
	req := &v1.MenuRequest{
		Name:      "Updated Menu",
		Path:      "/updated",
		Component: "@/pages/Updated",
		Icon:      "edit",
	}

	mockMenuRepo.EXPECT().Update(ctx, menuId, gomock.Any()).Return(nil)

	err := menuService.Update(ctx, menuId, req)

	assert.NoError(t, err)
}

func TestMenuService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMenuRepo := mock_repository.NewMockMenuRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	menuService := service.NewMenuService(srv, mockMenuRepo)

	ctx := context.Background()
	menuId := uint(1)

	mockMenuRepo.EXPECT().Delete(ctx, menuId).Return(nil)

	err := menuService.Delete(ctx, menuId)

	assert.NoError(t, err)
}
