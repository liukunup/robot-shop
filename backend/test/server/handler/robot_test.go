package handler

import (
	v1 "backend/api/v1"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/model"
	"backend/test/mocks/service"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"gorm.io/gorm"
)

func TestRobotHandler_ListRobots(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRobotService := mock_service.NewMockRobotService(ctrl)
	mockRobotService.EXPECT().List(gomock.Any(), gomock.Any()).Return(&v1.RobotSearchResponseData{
		Total: 2,
		List: []v1.RobotDataItem{
			{
				Id:       1,
				Name:     "Test Robot 1",
				Owner:    userId,
				Desc:     "Test Description",
				Webhook:  "https://example.com/webhook1",
				Enabled:  true,
			},
			{
				Id:       2,
				Name:     "Test Robot 2",
				Owner:    userId,
				Desc:     "Test Description 2",
				Webhook:  "https://example.com/webhook2",
				Enabled:  true,
			},
		},
	}, nil)

	robotHandler := handler.NewRobotHandler(hdl, mockRobotService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.GET("/robots", robotHandler.ListRobots)

	e := newHttpExcept(t, router)
	obj := e.GET("/robots").
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

func TestRobotHandler_CreateRobot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.RobotRequest{
		Name:    "New Robot",
		Desc:    "New Description",
		Webhook: "https://example.com/webhook",
		Enabled: true,
		Owner:   userId,
	}

	mockRobotService := mock_service.NewMockRobotService(ctrl)
	mockRobotService.EXPECT().Create(gomock.Any(), &params).Return(nil)

	robotHandler := handler.NewRobotHandler(hdl, mockRobotService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.POST("/robots", robotHandler.CreateRobot)

	e := newHttpExcept(t, router)
	obj := e.POST("/robots").
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

func TestRobotHandler_UpdateRobot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	robotId := uint(1)
	params := v1.RobotRequest{
		Name:    "Updated Robot",
		Desc:    "Updated Description",
		Webhook: "https://example.com/webhook-updated",
		Enabled: true,
		Owner:   userId,
	}

	mockRobotService := mock_service.NewMockRobotService(ctrl)
	mockRobotService.EXPECT().Update(gomock.Any(), robotId, &params).Return(nil)

	robotHandler := handler.NewRobotHandler(hdl, mockRobotService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.PUT("/robots/:id", robotHandler.UpdateRobot)

	e := newHttpExcept(t, router)
	obj := e.PUT("/robots/1").
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

func TestRobotHandler_GetRobot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	robotId := uint(1)
	mockRobotService := mock_service.NewMockRobotService(ctrl)
	mockRobotService.EXPECT().Get(gomock.Any(), robotId).Return(model.Robot{
		Model: gorm.Model{
			ID: 1,
		},
		Name:    "Test Robot",
		Owner:   userId,
		Desc:    "Test Description",
		Webhook: "https://example.com/webhook",
		Enabled: true,
	}, nil)

	robotHandler := handler.NewRobotHandler(hdl, mockRobotService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.GET("/robots/:id", robotHandler.GetRobot)

	e := newHttpExcept(t, router)
	obj := e.GET("/robots/1").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("success").IsEqual(true)
	obj.Value("errorMessage").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("id").IsEqual(1)
	objData.Value("name").IsEqual("Test Robot")
}

func TestRobotHandler_DeleteRobot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	robotId := uint(1)
	mockRobotService := mock_service.NewMockRobotService(ctrl)
	mockRobotService.EXPECT().Delete(gomock.Any(), robotId).Return(nil)

	robotHandler := handler.NewRobotHandler(hdl, mockRobotService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.DELETE("/robots/:id", robotHandler.DeleteRobot)

	e := newHttpExcept(t, router)
	obj := e.DELETE("/robots/1").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("success").IsEqual(true)
	obj.Value("errorMessage").IsEqual("ok")
}

