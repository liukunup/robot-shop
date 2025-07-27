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
// @Param page query int false "页码"
// @Param pageSize query int false "分页大小"
// @Param name query string false "角色名"
// @Param casbinRole query string false "Casbin Role"
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
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, data)
}

// CreateRole godoc
// @Summary 创建角色
// @Schemes
// @Description 创建一个新的角色
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.RoleRequest true "角色数据"
// @Success 200 {object} v1.Response
// @Router /admin/roles [post]
// @ID CreateRole
func (h *RoleHandler) CreateRole(ctx *gin.Context) {
	var req v1.RoleRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("CreateRole bind error", zap.Error(err))
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

// UpdateRole godoc
// @Summary 更新角色
// @Schemes
// @Description 目前只允许更新角色名称
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path uint true "角色ID"
// @Param request body v1.RoleRequest true "角色数据"
// @Success 200 {object} v1.Response
// @Router /admin/roles/{id} [put]
// @ID UpdateRole
func (h *RoleHandler) UpdateRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("UpdateRole parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req v1.RoleRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("UpdateRole bind error", zap.Error(err))
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

// DeleteRole godoc
// @Summary 删除角色
// @Schemes
// @Description 删除指定ID的角色
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path uint true "角色ID"
// @Success 200 {object} v1.Response
// @Router /admin/roles/{id} [delete]
// @ID DeleteRole
func (h *RoleHandler) DeleteRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("DeleteRole parse id error", zap.Error(err))
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

// GetRolePermissions godoc
// @Summary 获取角色权限
// @Schemes
// @Description 获取指定角色的权限列表
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param role query string true "角色名"
// @Success 200 {object} v1.GetRolePermissionResponse
// @Router /admin/roles/permissions [get]
// @ID GetRolePermissions
func (h *RoleHandler) GetRolePermissions(ctx *gin.Context) {
	var req v1.GetRolePermissionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("GetRolePermissions bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	data, err := h.roleService.GetPermissions(ctx, req.CasbinRole)
	if err != nil {
		h.logger.WithContext(ctx).Error("roleService.GetPermissions error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, data)
}

// UpdateRolePermissions godoc
// @Summary 更新角色权限
// @Schemes
// @Description 更新指定角色的权限列表
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateRolePermissionRequest true "更新参数"
// @Success 200 {object} v1.Response
// @Router /admin/roles/permissions [put]
// @ID UpdateRolePermissions
func (h *RoleHandler) UpdateRolePermissions(ctx *gin.Context) {
	var req v1.UpdateRolePermissionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("UpdateRolePermissions bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	err := h.roleService.UpdatePermissions(ctx, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("roleService.UpdatePermissions error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}
