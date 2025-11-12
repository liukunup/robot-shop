package handler

import (
	v1 "backend/api/v1"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/test/mocks/service"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestMenuHandler_ListMenus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMenuService := mock_service.NewMockMenuService(ctrl)
	mockMenuService.EXPECT().List(gomock.Any(), gomock.Any()).Return(&v1.MenuSearchResponseData{
		Total: 2,
		List: []v1.MenuDataItem{
			{
				ID:        1,
				Name:      "Dashboard",
				Path:      "/dashboard",
				Component: "@/pages/Dashboard",
				Icon:      "dashboard",
			},
			{
				ID:        2,
				Name:      "Admin",
				Path:      "/admin",
				Component: "@/pages/Admin",
				Icon:      "crown",
			},
		},
	}, nil)

	menuHandler := handler.NewMenuHandler(hdl, mockMenuService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.GET("/admin/menus", menuHandler.ListMenus)

	e := newHttpExcept(t, router)
	obj := e.GET("/admin/menus").
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

func TestMenuHandler_CreateMenu(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.MenuRequest{
		Name:      "New Menu",
		Path:      "/new",
		Component: "@/pages/New",
		Icon:      "star",
	}

	mockMenuService := mock_service.NewMockMenuService(ctrl)
	mockMenuService.EXPECT().Create(gomock.Any(), &params).Return(nil)

	menuHandler := handler.NewMenuHandler(hdl, mockMenuService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.POST("/admin/menus", menuHandler.CreateMenu)

	e := newHttpExcept(t, router)
	obj := e.POST("/admin/menus").
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

func TestMenuHandler_UpdateMenu(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	menuId := uint(1)
	params := v1.MenuRequest{
		Name:      "Updated Menu",
		Path:      "/updated",
		Component: "@/pages/Updated",
		Icon:      "edit",
	}

	mockMenuService := mock_service.NewMockMenuService(ctrl)
	mockMenuService.EXPECT().Update(gomock.Any(), menuId, &params).Return(nil)

	menuHandler := handler.NewMenuHandler(hdl, mockMenuService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.PUT("/admin/menus/:id", menuHandler.UpdateMenu)

	e := newHttpExcept(t, router)
	obj := e.PUT("/admin/menus/1").
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

func TestMenuHandler_DeleteMenu(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	menuId := uint(1)
	mockMenuService := mock_service.NewMockMenuService(ctrl)
	mockMenuService.EXPECT().Delete(gomock.Any(), menuId).Return(nil)

	menuHandler := handler.NewMenuHandler(hdl, mockMenuService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.DELETE("/admin/menus/:id", menuHandler.DeleteMenu)

	e := newHttpExcept(t, router)
	obj := e.DELETE("/admin/menus/1").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("success").IsEqual(true)
	obj.Value("errorMessage").IsEqual("ok")
}
