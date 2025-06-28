package handler

import (
	v1 "backend/api/v1"
	"backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RoleHandler struct {
	*Handler
	roleService service.RoleService
}

func NewRoleHandler(
	handler *Handler,
	roleService service.RoleService,
) *RoleHandler {
	return &RoleHandler{
		Handler:     handler,
		roleService: roleService,
	}
}

// ListRoles godoc
// @Summary 获取角色列表
// @Schemes
// @Description 搜索时支持角色名和 Casbin Role 筛选
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "页码"
// @Param pageSize query int true "分页大小"
// @Param name query string false "角色名"
// @Param role query string false "Casbin Role"
// @Success 200 {object} v1.RoleSearchResponse
// @Router /admin/roles [get]
// @ID ListRoles
func (h *RoleHandler) ListRoles(ctx *gin.Context) {
	var req v1.RoleSearchRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("ListRoles bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	data, err := h.roleService.List(ctx, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("roleService.List error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, data)
}

// RoleCreate godoc
// @Summary 创建角色
// @Schemes
// @Description 创建一个新的角色
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.RoleRequest true "角色信息"
// @Success 200 {object} v1.Response
// @Router /admin/roles [post]
// @ID RoleCreate
func (h *RoleHandler) RoleCreate(ctx *gin.Context) {
	var req v1.RoleRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("RoleCreate bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.roleService.Create(ctx, &req); err != nil {
		h.logger.WithContext(ctx).Error("roleService.Create error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// RoleUpdate godoc
// @Summary 更新角色
// @Schemes
// @Description 目前只允许更新角色名称
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path uint true "角色ID"
// @Param request body v1.RoleRequest true "角色信息"
// @Success 200 {object} v1.Response
// @Router /admin/roles/{id} [put]
// @ID RoleUpdate
func (h *RoleHandler) RoleUpdate(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("RoleUpdate parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req v1.RoleRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("RoleUpdate bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.roleService.Update(ctx, uint(id), &req); err != nil {
		h.logger.WithContext(ctx).Error("roleService.Update error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// RoleDelete godoc
// @Summary 删除角色
// @Schemes
// @Description
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path uint true "角色ID"
// @Success 200 {object} v1.Response
// @Router /admin/roles/{id} [delete]
// @ID RoleDelete
func (h *RoleHandler) RoleDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("RoleDelete parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.roleService.Delete(ctx, uint(id)); err != nil {
		h.logger.WithContext(ctx).Error("roleService.Delete error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// GetRolePermission godoc
// @Summary 获取角色权限
// @Schemes
// @Description 获取指定角色的权限列表
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param role query string true "角色名"
// @Success 200 {object} v1.GetRolePermissionResponse
// @Router /admin/roles/permission [get]
// @ID GetRolePermission
func (h *RoleHandler) GetRolePermission(ctx *gin.Context) {
	var req v1.GetRolePermissionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("GetRolePermission bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	data, err := h.roleService.GetRolePermission(ctx, req.Role)
	if err != nil {
		h.logger.WithContext(ctx).Error("roleService.GetRolePermission error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, data)
}

// UpdateRolePermission godoc
// @Summary 更新角色权限
// @Schemes
// @Description 更新指定角色的权限列表
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateRolePermissionRequest true "更新参数"
// @Success 200 {object} v1.Response
// @Router /admin/roles/permission [put]
// @ID UpdateRolePermission
func (h *RoleHandler) UpdateRolePermission(ctx *gin.Context) {
	var req v1.UpdateRolePermissionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("UpdateRolePermission bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	err := h.roleService.UpdateRolePermission(ctx, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("roleService.UpdateRolePermission error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}
