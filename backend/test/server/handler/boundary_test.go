package handler

import (
	v1 "backend/api/v1"
	"backend/internal/handler"
	"backend/internal/middleware"
	mock_service "backend/test/mocks/service"
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

// User模块边界条件测试

func TestUserHandler_Register_LongEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 测试超长邮箱
	longEmail := ""
	for i := 0; i < 100; i++ {
		longEmail += "a"
	}
	longEmail += "@example.com"

	params := v1.RegisterRequest{
		Email:    longEmail,
		Password: "123456",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().Register(gomock.Any(), &params).Return(nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)

	testRouter := gin.New()
	testRouter.Use(middleware.CORSMiddleware())
	testRouter.POST("/register", userHandler.Register)

	e := newHttpExcept(t, testRouter)
	obj := e.POST("/register").
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("success").IsEqual(true)
}

func TestUserHandler_Register_ShortPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.RegisterRequest{
		Email:    "test@example.com",
		Password: "12", // 过短的密码
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().Register(gomock.Any(), &params).Return(errors.New("password too short"))

	userHandler := handler.NewUserHandler(hdl, mockUserService)

	testRouter := gin.New()
	testRouter.Use(middleware.CORSMiddleware())
	testRouter.POST("/register", userHandler.Register)

	e := newHttpExcept(t, testRouter)
	e.POST("/register").
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusInternalServerError)
}

func TestUserHandler_Login_EmptyUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.LoginRequest{
		Username: "", // 空用户名
		Password: "123456",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	userHandler := handler.NewUserHandler(hdl, mockUserService)

	testRouter := gin.New()
	testRouter.Use(middleware.CORSMiddleware())
	testRouter.POST("/login", userHandler.Login)

	e := newHttpExcept(t, testRouter)
	e.POST("/login").
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusBadRequest) // binding validation会失败
}

func TestUserHandler_Login_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().Login(gomock.Any(), &params).Return(nil, v1.ErrUnauthorized)

	userHandler := handler.NewUserHandler(hdl, mockUserService)

	testRouter := gin.New()
	testRouter.Use(middleware.CORSMiddleware())
	testRouter.POST("/login", userHandler.Login)

	e := newHttpExcept(t, testRouter)
	e.POST("/login").
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusUnauthorized)
}

// Robot模块边界条件测试

func TestRobotHandler_CreateRobot_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.RobotRequest{
		Name:    "", // 空名称
		Desc:    "Test",
		Webhook: "https://example.com/webhook",
		Enabled: true,
		Owner:   "test",
	}

	mockRobotService := mock_service.NewMockRobotService(ctrl)
	mockRobotService.EXPECT().Create(gomock.Any(), &params).Return(errors.New("name cannot be empty"))

	robotHandler := handler.NewRobotHandler(hdl, mockRobotService)

	testRouter := gin.New()
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.POST("/robots", robotHandler.CreateRobot)

	e := newHttpExcept(t, testRouter)
	e.POST("/robots").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusInternalServerError)
}

func TestRobotHandler_CreateRobot_InvalidWebhook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.RobotRequest{
		Name:    "Test Robot",
		Desc:    "Test",
		Webhook: "not-a-valid-url", // 无效的URL
		Enabled: true,
		Owner:   "test",
	}

	mockRobotService := mock_service.NewMockRobotService(ctrl)
	mockRobotService.EXPECT().Create(gomock.Any(), &params).Return(nil) // 业务层可能不校验

	robotHandler := handler.NewRobotHandler(hdl, mockRobotService)

	testRouter := gin.New()
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.POST("/robots", robotHandler.CreateRobot)

	e := newHttpExcept(t, testRouter)
	obj := e.POST("/robots").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("success").IsEqual(true)
}

// Role模块边界条件测试

func TestRoleHandler_CreateRole_DuplicateCasbinRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.RoleRequest{
		Name:       "Admin",
		CasbinRole: "admin", // 重复的角色
	}

	mockRoleService := mock_service.NewMockRoleService(ctrl)
	mockRoleService.EXPECT().Create(gomock.Any(), &params).Return(errors.New("duplicate casbin role"))

	roleHandler := handler.NewRoleHandler(hdl, mockRoleService)

	testRouter := gin.New()
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.POST("/admin/roles", roleHandler.CreateRole)

	e := newHttpExcept(t, testRouter)
	e.POST("/admin/roles").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusInternalServerError)
}

func TestRoleHandler_DeleteRole_WithUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleService := mock_service.NewMockRoleService(ctrl)
	mockRoleService.EXPECT().Delete(gomock.Any(), uint(1)).Return(errors.New("role is in use"))

	roleHandler := handler.NewRoleHandler(hdl, mockRoleService)

	testRouter := gin.New()
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.DELETE("/admin/roles/:id", roleHandler.DeleteRole)

	e := newHttpExcept(t, testRouter)
	e.DELETE("/admin/roles/1").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusInternalServerError)
}

// Menu模块边界条件测试

func TestMenuHandler_CreateMenu_InvalidParentId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.MenuRequest{
		ParentID:  999, // 不存在的父菜单ID
		Icon:      "test",
		Name:      "Test Menu",
		Path:      "/test",
		Component: "@/pages/Test",
	}

	mockMenuService := mock_service.NewMockMenuService(ctrl)
	mockMenuService.EXPECT().Create(gomock.Any(), &params).Return(errors.New("parent menu not found"))

	menuHandler := handler.NewMenuHandler(hdl, mockMenuService)

	testRouter := gin.New()
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.POST("/admin/menus", menuHandler.CreateMenu)

	e := newHttpExcept(t, testRouter)
	e.POST("/admin/menus").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusInternalServerError)
}

func TestMenuHandler_UpdateMenu_CircularReference(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.MenuRequest{
		ParentID:  1, // 设置自己为父菜单造成循环引用
		Icon:      "test",
		Name:      "Test Menu",
		Path:      "/test",
		Component: "@/pages/Test",
	}

	mockMenuService := mock_service.NewMockMenuService(ctrl)
	mockMenuService.EXPECT().Update(gomock.Any(), uint(1), &params).Return(errors.New("circular reference detected"))

	menuHandler := handler.NewMenuHandler(hdl, mockMenuService)

	testRouter := gin.New()
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.PUT("/admin/menus/:id", menuHandler.UpdateMenu)

	e := newHttpExcept(t, testRouter)
	e.PUT("/admin/menus/1").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusInternalServerError)
}

// API模块边界条件测试

func TestApiHandler_CreateApi_DuplicatePath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.ApiRequest{
		Group:  "User",
		Name:   "ListUsers",
		Path:   "/v1/admin/users", // 重复的路径
		Method: "GET",
	}

	mockApiService := mock_service.NewMockApiService(ctrl)
	mockApiService.EXPECT().Create(gomock.Any(), &params).Return(errors.New("duplicate api path"))

	apiHandler := handler.NewApiHandler(hdl, mockApiService)

	testRouter := gin.New()
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.POST("/admin/apis", apiHandler.CreateApi)

	e := newHttpExcept(t, testRouter)
	e.POST("/admin/apis").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusInternalServerError)
}

func TestApiHandler_ListApis_LargePageSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApiService := mock_service.NewMockApiService(ctrl)
	apiHandler := handler.NewApiHandler(hdl, mockApiService)

	testRouter := gin.New()
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.GET("/admin/apis", apiHandler.ListApis)

	e := newHttpExcept(t, testRouter)
	e.GET("/admin/apis").
		WithQuery("page", 1).
		WithQuery("pageSize", 1000). // 超大的pageSize
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusBadRequest) // 应该被binding validation拦截
}
