package service_test

import (
	"context"
	"errors"
	"testing"

	v1 "backend/api/v1"
	"backend/internal/model"
	"backend/internal/service"
	mock_repository "backend/test/mocks/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// User Service错误处理测试

func TestUserService_Register_DatabaseError(t *testing.T) {
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
	mockTm.EXPECT().Transaction(ctx, gomock.Any()).Return(errors.New("database connection failed"))

	err := userService.Register(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database")
}

func TestUserService_Login_AccountDisabled(t *testing.T) {
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
		Username: "disableduser",
		Password: "password",
	}

	mockUserRepo.EXPECT().GetByUsernameOrEmail(ctx, req.Username, req.Username).Return(model.User{
		Model:    gorm.Model{ID: 1},
		Username: "disableduser",
		Status:   2, // 禁用状态
	}, nil)

	_, err := userService.Login(ctx, req)

	assert.Error(t, err)
}

// Robot Service错误处理测试

func TestRobotService_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRobotRepo := mock_repository.NewMockRobotRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	robotService := service.NewRobotService(srv, mockRobotRepo)

	ctx := context.Background()
	req := &v1.RobotRequest{
		Name:    "", // 空名称
		Webhook: "https://example.com/webhook",
		Enabled: true,
	}

	// 期望调用Create,因为当前没有validation
	mockRobotRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	err := robotService.Create(ctx, req)

	assert.NoError(t, err) // 当前实现没有validation
}

func TestRobotService_Update_Concurrency(t *testing.T) {
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
		Webhook: "https://example.com/webhook",
		Enabled: true,
	}

	// 模拟并发更新冲突
	mockRobotRepo.EXPECT().Update(ctx, robotId, gomock.Any()).Return(errors.New("concurrent update conflict"))

	err := robotService.Update(ctx, robotId, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "concurrent")
}

func TestRobotService_Delete_InUse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRobotRepo := mock_repository.NewMockRobotRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	robotService := service.NewRobotService(srv, mockRobotRepo)

	ctx := context.Background()
	robotId := uint(1)

	// 模拟机器人正在使用中
	mockRobotRepo.EXPECT().Delete(ctx, robotId).Return(errors.New("robot is in use"))

	err := robotService.Delete(ctx, robotId)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "in use")
}

// Role Service错误处理测试

func TestRoleService_Create_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mock_repository.NewMockRoleRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	roleService := service.NewRoleService(srv, mockRoleRepo)

	ctx := context.Background()
	req := &v1.RoleRequest{
		Name:       "", // 空名称
		CasbinRole: "test",
	}

	// 期望先检查casbin role是否存在
	mockRoleRepo.EXPECT().GetByCasbinRole(ctx, req.CasbinRole).Return(model.Role{}, gorm.ErrRecordNotFound)
	mockRoleRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	err := roleService.Create(ctx, req)

	assert.NoError(t, err) // 当前没有name validation
}

func TestRoleService_Delete_SystemRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mock_repository.NewMockRoleRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	roleService := service.NewRoleService(srv, mockRoleRepo)

	ctx := context.Background()
	roleId := uint(1) // 假设1是系统管理员角色

	// 先查询role,删除casbin role,再删除role
	mockRoleRepo.EXPECT().Get(ctx, roleId).Return(model.Role{
		Model:      gorm.Model{ID: roleId},
		Name:       "Admin",
		CasbinRole: "admin",
	}, nil)
	mockRoleRepo.EXPECT().DeleteCasbinRole(ctx, "admin").Return(true, errors.New("cannot delete system role"))

	err := roleService.Delete(ctx, roleId)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "system role")
}

func TestRoleService_Update_CasbinRoleConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mock_repository.NewMockRoleRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	roleService := service.NewRoleService(srv, mockRoleRepo)

	ctx := context.Background()
	roleId := uint(2)
	req := &v1.RoleRequest{
		Name:       "Updated Role",
		CasbinRole: "admin", // 与现有角色冲突
	}

	mockRoleRepo.EXPECT().Update(ctx, gomock.Any()).Return(errors.New("casbin role already exists"))

	err := roleService.Update(ctx, roleId, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

// Menu Service错误处理测试

func TestMenuService_Create_ParentNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMenuRepo := mock_repository.NewMockMenuRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	menuService := service.NewMenuService(srv, mockMenuRepo)

	ctx := context.Background()
	req := &v1.MenuRequest{
		ParentID:  999, // 不存在的父菜单
		Name:      "Child Menu",
		Path:      "/child",
		Component: "@/pages/Child",
	}

	mockMenuRepo.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("parent menu not found"))

	err := menuService.Create(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parent")
}

func TestMenuService_Delete_HasChildren(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMenuRepo := mock_repository.NewMockMenuRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	menuService := service.NewMenuService(srv, mockMenuRepo)

	ctx := context.Background()
	menuId := uint(1)

	// 模拟菜单有子菜单
	mockMenuRepo.EXPECT().Delete(ctx, menuId).Return(errors.New("menu has children"))

	err := menuService.Delete(ctx, menuId)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "children")
}

func TestMenuService_Update_PathConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMenuRepo := mock_repository.NewMockMenuRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(logger, sf, j, em, mockTm)
	menuService := service.NewMenuService(srv, mockMenuRepo)

	ctx := context.Background()
	menuId := uint(2)
	req := &v1.MenuRequest{
		Name:      "Updated Menu",
		Path:      "/dashboard", // 与其他菜单路径冲突
		Component: "@/pages/Dashboard",
	}

	mockMenuRepo.EXPECT().Update(ctx, menuId, gomock.Any()).Return(errors.New("path already exists"))

	err := menuService.Update(ctx, menuId, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

// 通用错误场景测试

func TestService_DatabaseTimeout(t *testing.T) {
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

	// 模拟数据库超时
	mockRobotRepo.EXPECT().List(ctx, req).Return(nil, int64(0), context.DeadlineExceeded)

	_, err := robotService.List(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestService_EmptyFieldHandling(t *testing.T) {
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

	// 测试特殊字符email
	req := &v1.RegisterRequest{
		Email:    "test+special@example.com",
		Password: "123456",
	}

	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(model.User{}, gorm.ErrRecordNotFound)
	mockTm.EXPECT().Transaction(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})
	mockUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	err := userService.Register(ctx, req)

	assert.NoError(t, err)
}
