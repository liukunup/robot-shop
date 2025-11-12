package handler

import (
	v1 "backend/api/v1"
	"backend/internal/handler"
	"backend/internal/middleware"
	mock_service "backend/test/mocks/service"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestUserHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.RegisterRequest{
		Email:    "test@example.com",
		Password: "123456",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().Register(gomock.Any(), &params).Return(nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)

	// Create a new router for this test to avoid JWT middleware
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
	obj.Value("errorMessage").IsEqual("ok")
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.LoginRequest{
		Username: "testuser",
		Password: "123456",
	}

	tokenPair := &v1.TokenPair{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.refresh",
		ExpiresIn:    900,
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().Login(gomock.Any(), &params).Return(tokenPair, nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)

	// Create a new router for this test to avoid JWT middleware
	testRouter := gin.New()
	testRouter.Use(middleware.CORSMiddleware())
	testRouter.POST("/login", userHandler.Login)

	obj := newHttpExcept(t, testRouter).POST("/login").
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("success").IsEqual(true)
	obj.Value("errorMessage").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("accessToken").IsEqual(tokenPair.AccessToken)
	objData.Value("refreshToken").IsEqual(tokenPair.RefreshToken)
}

func TestUserHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userId := uint(1)
	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().Get(gomock.Any(), userId).Return(&v1.UserDataItem{
		ID:       userId,
		Username: "testuser",
		Nickname: "Test User",
		Email:    "test@example.com",
	}, nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.GET("/users/:id", userHandler.GetUserByID)

	obj := newHttpExcept(t, router).GET("/users/1").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("success").IsEqual(true)
	obj.Value("errorMessage").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("username").IsEqual("testuser")
	objData.Value("nickname").IsEqual("Test User")
}

func TestUserHandler_UpdatePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.UpdatePasswordRequest{
		OldPassword: "oldpass123",
		NewPassword: "newpass456",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().UpdatePassword(gomock.Any(), uint(1), &params).Return(nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.PUT("/user/password", userHandler.UpdatePassword)

	obj := newHttpExcept(t, router).PUT("/user/password").
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("success").IsEqual(true)
	obj.Value("errorMessage").IsEqual("ok")
}

func TestUserHandler_ListUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().List(gomock.Any(), gomock.Any()).Return(&v1.UserSearchResponseData{
		Total: 2,
		List: []v1.UserDataItem{
			{
				ID:       1,
				Username: "user1",
				Nickname: "User One",
				Email:    "user1@example.com",
			},
			{
				ID:       2,
				Username: "user2",
				Nickname: "User Two",
				Email:    "user2@example.com",
			},
		},
	}, nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.GET("/admin/users", userHandler.ListUsers)

	obj := newHttpExcept(t, router).GET("/admin/users").
		WithQuery("page", 1).
		WithQuery("pageSize", 10).
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("success").IsEqual(true)
	obj.Value("errorMessage").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("total").IsEqual(2)
	objData.Value("list").Array().Length().IsEqual(2)
}
