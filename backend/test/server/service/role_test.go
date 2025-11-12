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

func TestRoleService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mock_repository.NewMockRoleRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	roleService := service.NewRoleService(srv, mockRoleRepo)

	ctx := context.Background()
	req := &v1.RoleSearchRequest{
		Page:     1,
		PageSize: 10,
	}

	mockRoles := []model.Role{
		{
			Model:      gorm.Model{ID: 1},
			Name:       "Admin",
			CasbinRole: "admin",
		},
		{
			Model:      gorm.Model{ID: 2},
			Name:       "User",
			CasbinRole: "user",
		},
	}

	mockRoleRepo.EXPECT().List(ctx, req).Return(mockRoles, int64(2), nil)

	result, err := roleService.List(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.Total)
	assert.Len(t, result.List, 2)
	assert.Equal(t, "Admin", result.List[0].Name)
}

func TestRoleService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mock_repository.NewMockRoleRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	roleService := service.NewRoleService(srv, mockRoleRepo)

	ctx := context.Background()
	req := &v1.RoleRequest{
		Name:       "New Role",
		CasbinRole: "newrole",
	}

	mockRoleRepo.EXPECT().GetByCasbinRole(ctx, req.CasbinRole).Return(model.Role{}, gorm.ErrRecordNotFound)
	mockRoleRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	err := roleService.Create(ctx, req)

	assert.NoError(t, err)
}

func TestRoleService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mock_repository.NewMockRoleRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	roleService := service.NewRoleService(srv, mockRoleRepo)

	ctx := context.Background()
	roleId := uint(1)
	req := &v1.RoleRequest{
		Name:       "Updated Role",
		CasbinRole: "updatedrole",
	}

	mockRoleRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	err := roleService.Update(ctx, roleId, req)

	assert.NoError(t, err)
}

func TestRoleService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mock_repository.NewMockRoleRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	roleService := service.NewRoleService(srv, mockRoleRepo)

	ctx := context.Background()
	roleId := uint(1)

	mockRoleRepo.EXPECT().Get(ctx, roleId).Return(model.Role{
		Model:      gorm.Model{ID: roleId},
		CasbinRole: "test",
	}, nil)
	mockRoleRepo.EXPECT().DeleteCasbinRole(ctx, "test").Return(true, nil)
	mockRoleRepo.EXPECT().Delete(ctx, roleId).Return(nil)

	err := roleService.Delete(ctx, roleId)

	assert.NoError(t, err)
}

func TestRoleService_ListAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mock_repository.NewMockRoleRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	roleService := service.NewRoleService(srv, mockRoleRepo)

	ctx := context.Background()

	mockRoles := []model.Role{
		{Model: gorm.Model{ID: 1}, Name: "Admin", CasbinRole: "admin"},
		{Model: gorm.Model{ID: 2}, Name: "User", CasbinRole: "user"},
		{Model: gorm.Model{ID: 3}, Name: "Guest", CasbinRole: "guest"},
	}

	mockRoleRepo.EXPECT().ListAll(ctx).Return(mockRoles, nil)

	result, err := roleService.ListAll(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 3)
}
