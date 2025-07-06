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
// @Param page query int true "页码"
// @Param pageSize query int true "分页大小"
// @Param name query string false "名称"
// @Param path query string false "路径"
// @Param access query string false "可见性"
// @Success 200 {object} v1.MenuSearchResponse
// @Router /admin/menus [get]
// @ID ListMenus
func (h *MenuHandler) ListMenus(ctx *gin.Context) {
	var req v1.MenuSearchRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("ListMenus bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	data, err := h.menuService.List(ctx, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("menuService.List error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, data)
}

// CreateMenu godoc
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
// @ID CreateMenu
func (h *MenuHandler) CreateMenu(ctx *gin.Context) {
	var req v1.MenuRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("CreateMenu bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.menuService.Create(ctx, &req); err != nil {
		h.logger.WithContext(ctx).Error("menuService.Create error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// UpdateMenu godoc
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
// @ID UpdateMenu
func (h *MenuHandler) UpdateMenu(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("UpdateMenu parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req v1.MenuRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("UpdateMenu bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.menuService.Update(ctx, uint(id), &req); err != nil {
		h.logger.WithContext(ctx).Error("menuService.Update error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// DeleteMenu godoc
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
// @ID DeleteMenu
func (h *MenuHandler) DeleteMenu(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("DeleteMenu parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.menuService.Delete(ctx, uint(id)); err != nil {
		h.logger.WithContext(ctx).Error("menuService.Delete error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return

	}
	v1.HandleSuccess(ctx, nil)
}
