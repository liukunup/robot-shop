package handler

import (
	v1 "backend/api/v1"
	"backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MenuHandler struct {
	*Handler
	menuService service.MenuService
}

func NewMenuHandler(
	handler *Handler,
	menuService service.MenuService,
) *MenuHandler {
	return &MenuHandler{
		Handler:     handler,
		menuService: menuService,
	}
}

// ListMenus godoc
// @Summary 获取菜单列表
// @Schemes
// @Description 获取所有菜单
// @Tags Menu
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.MenuSearchResponse
// @Router /admin/menus [get]
// @ID ListMenus
func (h *MenuHandler) ListMenus(ctx *gin.Context) {
	data, err := h.menuService.ListMenus(ctx)
	if err != nil {
		h.logger.WithContext(ctx).Error("menuService.ListMenus error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, data)
}

// MenuCreate godoc
// @Summary 创建菜单
// @Schemes
// @Description 创建一个新的菜单
// @Tags Menu
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.MenuRequest true "菜单数据"
// @Success 200 {object} v1.Response
// @Router /admin/menus [post]
// @ID MenuCreate
func (h *MenuHandler) MenuCreate(ctx *gin.Context) {
	var req v1.MenuRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("MenuCreate bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.menuService.MenuCreate(ctx, &req); err != nil {
		h.logger.WithContext(ctx).Error("menuService.MenuCreate error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// MenuUpdate godoc
// @Summary 更新菜单
// @Schemes
// @Description 更新菜单数据
// @Tags Menu
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path uint true "菜单ID"
// @Param request body v1.MenuRequest true "菜单数据"
// @Success 200 {object} v1.Response
// @Router /admin/menus/{id} [put]
// @ID MenuUpdate
func (h *MenuHandler) MenuUpdate(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("MenuUpdate parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req v1.MenuRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("MenuUpdate bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.menuService.MenuUpdate(ctx, uint(id), &req); err != nil {
		h.logger.WithContext(ctx).Error("menuService.MenuUpdate error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// MenuDelete godoc
// @Summary 删除菜单
// @Schemes
// @Description 删除指定ID的菜单
// @Tags Menu
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path uint true "菜单ID"
// @Success 200 {object} v1.Response
// @Router /admin/menus/{id} [delete]
// @ID MenuDelete
func (h *MenuHandler) MenuDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("MenuDelete parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.menuService.MenuDelete(ctx, uint(id)); err != nil {
		h.logger.WithContext(ctx).Error("menuService.MenuDelete error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return

	}
	v1.HandleSuccess(ctx, nil)
}
