package handler

import (
	v1 "backend/api/v1"
	"backend/internal/service"
	"backend/pkg/time"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ApiHandler struct {
	*Handler
	apiService service.ApiService
}

func NewApiHandler(
	handler *Handler,
	apiService service.ApiService,
) *ApiHandler {
	return &ApiHandler{
		Handler:    handler,
		apiService: apiService,
	}
}

// ListApis godoc
// @Summary 获取 API 列表
// @Schemes
// @Description 搜索时支持分组、名称、路径和方法字段的筛选
// @Tags API
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "页码"
// @Param pageSize query int true "分页大小"
// @Param group query string false "分组"
// @Param name query string false "名称"
// @Param path query string false "路径"
// @Param method query string false "方法"
// @Success 200 {object} v1.ApiSearchResponse
// @Router /admin/apis [get]
// @ID ListApis
func (h *ApiHandler) ListApis(ctx *gin.Context) {
	var req v1.ApiSearchRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("ListApis bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	data, err := h.apiService.List(ctx, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("apiService.List error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, data)
}

// CreateApi godoc
// @Summary 创建 API
// @Schemes
// @Description 实际上是录入一个新的 API
// @Tags API
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.ApiRequest true "API 数据"
// @Success 200 {object} v1.Response
// @Router /admin/apis [post]
// @ID CreateApi
func (h *ApiHandler) CreateApi(ctx *gin.Context) {
	var req v1.ApiRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("CreateApi bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.apiService.Create(ctx, &req); err != nil {
		h.logger.WithContext(ctx).Error("apiService.Create error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// UpdateApi godoc
// @Summary 更新 API
// @Schemes
// @Description 更新指定`ID`的 API
// @Tags API
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "ID"
// @Param request body v1.ApiRequest true "API 数据"
// @Success 200 {object} v1.Response
// @Router /admin/apis/{id} [put]
// @ID UpdateApi
func (h *ApiHandler) UpdateApi(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("UpdateApi parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req v1.ApiRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("UpdateApi bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.apiService.Update(ctx, uint(id), &req); err != nil {
		h.logger.WithContext(ctx).Error("apiService.Update error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// DeleteApi godoc
// @Summary 删除 API
// @Schemes
// @Description 删除指定`ID`的 API
// @Tags API
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query uint true "ID"
// @Success 200 {object} v1.Response
// @Router /admin/apis/{id} [delete]
// @ID DeleteApi
func (h *ApiHandler) DeleteApi(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("DeleteApi parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.apiService.Delete(ctx, uint(id)); err != nil {
		h.logger.WithContext(ctx).Error("apiService.Delete error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// GetApi godoc
// @Summary 获取 API 详情
// @Schemes
// @Description 获取指定`ID`的 API 数据
// @Tags API
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v1.ApiResponse
// @Router /admin/apis/{id} [get]
// @ID GetApi
func (h *ApiHandler) GetApi(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("GetApi parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	api, err := h.apiService.Get(ctx, uint(id))
	if err != nil {
		h.logger.WithContext(ctx).Error("apiService.Get error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.ApiDataItem{
		ID:        api.ID,
		CreatedAt: time.FormatTime(api.CreatedAt),
		UpdatedAt: time.FormatTime(api.UpdatedAt),
		Group:     api.Group,
		Name:      api.Name,
		Path:      api.Path,
		Method:    api.Method,
	})
}
