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
	"gorm.io/gorm"
)

func TestApiHandler_ListApis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApiService := mock_service.NewMockApiService(ctrl)
	mockApiService.EXPECT().List(gomock.Any(), gomock.Any()).Return(&v1.ApiSearchResponseData{
		Total: 2,
		List: []v1.ApiDataItem{
			{
				ID:     1,
				Group:  "User",
				Name:   "ListUsers",
				Path:   "/v1/admin/users",
				Method: "GET",
			},
			{
				ID:     2,
				Group:  "Robot",
				Name:   "ListRobots",
				Path:   "/v1/robots",
				Method: "GET",
			},
		},
	}, nil)

	apiHandler := handler.NewApiHandler(hdl, mockApiService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.GET("/admin/apis", apiHandler.ListApis)

	e := newHttpExcept(t, router)
	obj := e.GET("/admin/apis").
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

func TestApiHandler_CreateApi(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.ApiRequest{
		Group:  "Test",
		Name:   "TestApi",
		Path:   "/v1/test",
		Method: "POST",
	}

	mockApiService := mock_service.NewMockApiService(ctrl)
	mockApiService.EXPECT().Create(gomock.Any(), &params).Return(nil)

	apiHandler := handler.NewApiHandler(hdl, mockApiService)

	testRouter := gin.New()
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.POST("/admin/apis", apiHandler.CreateApi)

	e := newHttpExcept(t, testRouter)
	obj := e.POST("/admin/apis").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("success").IsEqual(true)
	obj.Value("errorMessage").IsEqual("ok")
}

func TestApiHandler_UpdateApi(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.ApiRequest{
		Group:  "Test",
		Name:   "UpdatedApi",
		Path:   "/v1/updated",
		Method: "PUT",
	}

	mockApiService := mock_service.NewMockApiService(ctrl)
	mockApiService.EXPECT().Update(gomock.Any(), uint(1), &params).Return(nil)

	apiHandler := handler.NewApiHandler(hdl, mockApiService)

	testRouter := gin.New()
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.PUT("/admin/apis/:id", apiHandler.UpdateApi)

	e := newHttpExcept(t, testRouter)
	obj := e.PUT("/admin/apis/1").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("success").IsEqual(true)
	obj.Value("errorMessage").IsEqual("ok")
}

func TestApiHandler_DeleteApi(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApiService := mock_service.NewMockApiService(ctrl)
	mockApiService.EXPECT().Delete(gomock.Any(), uint(1)).Return(nil)

	apiHandler := handler.NewApiHandler(hdl, mockApiService)

	testRouter := gin.New()
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.DELETE("/admin/apis/:id", apiHandler.DeleteApi)

	e := newHttpExcept(t, testRouter)
	obj := e.DELETE("/admin/apis/1").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("success").IsEqual(true)
	obj.Value("errorMessage").IsEqual("ok")
}

// 边界条件测试
func TestApiHandler_ListApis_InvalidPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApiService := mock_service.NewMockApiService(ctrl)
	apiHandler := handler.NewApiHandler(hdl, mockApiService)

	testRouter := gin.New()
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.GET("/admin/apis", apiHandler.ListApis)

	e := newHttpExcept(t, testRouter)
	e.GET("/admin/apis").
		WithQuery("page", 0). // 无效的page
		WithQuery("pageSize", 10).
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusBadRequest)
}

func TestApiHandler_CreateApi_EmptyRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.ApiRequest{} // 空请求

	mockApiService := mock_service.NewMockApiService(ctrl)
	mockApiService.EXPECT().Create(gomock.Any(), &params).Return(nil)

	apiHandler := handler.NewApiHandler(hdl, mockApiService)

	testRouter := gin.New()
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.POST("/admin/apis", apiHandler.CreateApi)

	e := newHttpExcept(t, testRouter)
	obj := e.POST("/admin/apis").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("success").IsEqual(true)
}

func TestApiHandler_UpdateApi_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.ApiRequest{
		Group:  "Test",
		Name:   "UpdatedApi",
		Path:   "/v1/updated",
		Method: "PUT",
	}

	mockApiService := mock_service.NewMockApiService(ctrl)
	mockApiService.EXPECT().Update(gomock.Any(), uint(999), &params).Return(gorm.ErrRecordNotFound)

	apiHandler := handler.NewApiHandler(hdl, mockApiService)

	testRouter := gin.New()
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.PUT("/admin/apis/:id", apiHandler.UpdateApi)

	e := newHttpExcept(t, testRouter)
	e.PUT("/admin/apis/999").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusInternalServerError)
}
