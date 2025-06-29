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
// @Summary 获取接口列表
// @Schemes
// @Description 搜索时支持分组名、名称、路径和方法筛选
// @Tags API
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "页码"
// @Param pageSize query int true "分页大小"
// @Param group query string false "分组名"
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

	data, err := h.apiService.ListApis(ctx, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("apiService.ListApis error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, data)
}

// ApiCreate godoc
// @Summary 创建接口
// @Schemes
// @Description 创建一个新的接口
// @Tags API
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.ApiRequest true "接口数据"
// @Success 200 {object} v1.Response
// @Router /admin/apis [post]
// @ID ApiCreate
func (h *ApiHandler) ApiCreate(ctx *gin.Context) {
	var req v1.ApiRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("ApiCreate bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.apiService.ApiCreate(ctx, &req); err != nil {
		h.logger.WithContext(ctx).Error("apiService.ApiCreate error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// ApiUpdate godoc
// @Summary 更新接口
// @Schemes
// @Description 更新接口数据
// @Tags API
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "接口ID"
// @Param request body v1.ApiRequest true "接口数据"
// @Success 200 {object} v1.Response
// @Router /admin/apis/{id} [put]
// @ID ApiUpdate
func (h *ApiHandler) ApiUpdate(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("ApiUpdate parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req v1.ApiRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("ApiUpdate bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.apiService.ApiUpdate(ctx, uint(id), &req); err != nil {
		h.logger.WithContext(ctx).Error("apiService.ApiUpdate error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// ApiDelete godoc
// @Summary 删除接口
// @Schemes
// @Description 删除指定ID的接口
// @Tags API
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query uint true "接口ID"
// @Success 200 {object} v1.Response
// @Router /admin/apis/{id} [delete]
// @ID ApiDelete
func (h *ApiHandler) ApiDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("ApiDelete parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.apiService.ApiDelete(ctx, uint(id)); err != nil {
		h.logger.WithContext(ctx).Error("apiService.ApiDelete error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// GetApi godoc
// @Summary 获取接口
// @Schemes
// @Description 获取指定ID的接口信息
// @Tags API
// @Accept json
// @Produce json
// @Param id path int true "接口ID"
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

	api, err := h.apiService.GetApi(ctx, uint(id))
	if err != nil {
		h.logger.WithContext(ctx).Error("apiService.GetApi error", zap.Error(err))
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
