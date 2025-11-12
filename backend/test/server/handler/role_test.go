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

func TestRoleHandler_ListRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleService := mock_service.NewMockRoleService(ctrl)
	mockRoleService.EXPECT().List(gomock.Any(), gomock.Any()).Return(&v1.RoleSearchResponseData{
		Total: 2,
		List: []v1.RoleDataItem{
			{
				ID:         1,
				Name:       "Admin",
				CasbinRole: "admin",
			},
			{
				ID:         2,
				Name:       "User",
				CasbinRole: "user",
			},
		},
	}, nil)

	roleHandler := handler.NewRoleHandler(hdl, mockRoleService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.GET("/admin/roles", roleHandler.ListRoles)

	e := newHttpExcept(t, router)
	obj := e.GET("/admin/roles").
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

func TestRoleHandler_CreateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.RoleRequest{
		Name:       "New Role",
		CasbinRole: "newrole",
	}

	mockRoleService := mock_service.NewMockRoleService(ctrl)
	mockRoleService.EXPECT().Create(gomock.Any(), &params).Return(nil)

	roleHandler := handler.NewRoleHandler(hdl, mockRoleService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.POST("/admin/roles", roleHandler.CreateRole)

	e := newHttpExcept(t, router)
	obj := e.POST("/admin/roles").
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

func TestRoleHandler_UpdateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleId := uint(1)
	params := v1.RoleRequest{
		Name:       "Updated Role",
		CasbinRole: "updatedrole",
	}

	mockRoleService := mock_service.NewMockRoleService(ctrl)
	mockRoleService.EXPECT().Update(gomock.Any(), roleId, &params).Return(nil)

	roleHandler := handler.NewRoleHandler(hdl, mockRoleService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.PUT("/admin/roles/:id", roleHandler.UpdateRole)

	e := newHttpExcept(t, router)
	obj := e.PUT("/admin/roles/1").
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

func TestRoleHandler_DeleteRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleId := uint(1)
	mockRoleService := mock_service.NewMockRoleService(ctrl)
	mockRoleService.EXPECT().Delete(gomock.Any(), roleId).Return(nil)

	roleHandler := handler.NewRoleHandler(hdl, mockRoleService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.DELETE("/admin/roles/:id", roleHandler.DeleteRole)

	e := newHttpExcept(t, router)
	obj := e.DELETE("/admin/roles/1").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("success").IsEqual(true)
	obj.Value("errorMessage").IsEqual("ok")
}
