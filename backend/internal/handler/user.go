package handler

import (
	v1 "backend/api/v1"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	*Handler
	userService service.UserService
}

func NewUserHandler(handler *Handler, userService service.UserService) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

// Register godoc
// @Summary 用户注册
// @Schemes
// @Description 目前只支持邮箱登录
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.RegisterRequest true "params"
// @Success 200 {object} v1.Response
// @Router /register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	req := new(v1.RegisterRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.Register(ctx, req); err != nil {
		h.logger.WithContext(ctx).Error("userService.Register error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// Login godoc
// @Summary 账号登录
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "params"
// @Success 200 {object} v1.LoginResponse
// @Router /login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	token, err := h.userService.Login(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.LoginResponseData{
		AccessToken: token,
	})
}

// GetAdminUsers godoc
// @Summary 获取管理员用户列表
// @Schemes
// @Description 获取管理员用户列表
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Param username query string false "用户名"
// @Param nickname query string false "昵称"
// @Param phone query string false "手机号"
// @Param email query string false "邮箱"
// @Success 200 {object} v1.GetAdminUsersResponse
// @Router /v1/admin/users [get]
func (h *UserHandler) GetUsers(ctx *gin.Context) {
	var req v1.GetUsersRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	data, err := h.userService.ListUsers(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, data)
}

// AdminUserUpdate godoc
// @Summary 更新管理员用户
// @Schemes
// @Description 更新管理员用户信息
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.AdminUserUpdateRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/user [put]
func (h *UserHandler) UserUpdate(ctx *gin.Context) {
	var req v1.AdminUserUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if err := h.userService.UserUpdate(ctx, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// AdminUserCreate godoc
// @Summary 创建管理员用户
// @Schemes
// @Description 创建新的管理员用户
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UserCreateRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/user [post]
func (h *UserHandler) UserCreate(ctx *gin.Context) {
	var req v1.AdminUserCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if err := h.userService.UserCreate(ctx, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// AdminUserDelete godoc
// @Summary 删除管理员用户
// @Schemes
// @Description 删除指定管理员用户
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query uint true "用户ID"
// @Success 200 {object} v1.Response
// @Router /v1/admin/user [delete]
func (h *UserHandler) AdminUserDelete(ctx *gin.Context) {
	var req v1.AdminUserDeleteRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if err := h.userService.UserDelete(ctx, req.ID); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return

	}
	v1.HandleSuccess(ctx, nil)
}

// GetAdminUser godoc
// @Summary 获取管理用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetAdminUserResponse
// @Router /v1/admin/user [get]
func (h *UserHandler) GetAdminUser(ctx *gin.Context) {
	data, err := h.userService.GetUser(ctx, GetUserIdFromCtx(ctx))
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, data)
}

// ListMenus godoc
// @Summary 获取用户菜单
// @Schemes
// @Description 获取当前用户的菜单列表
// @Tags 菜单模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetMenuResponse
// @Router /v1/menus [get]
func (h *UserHandler) ListMenus(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == 0 {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	data, err := h.userService.ListMenus(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, data)
}

// MenuUpdate godoc
// @Summary 更新菜单
// @Schemes
// @Description 更新菜单信息
// @Tags 菜单模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.MenuUpdateRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/menu [put]
func (h *UserHandler) MenuUpdate(ctx *gin.Context) {
	var req v1.MenuUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if err := h.userService.MenuUpdate(ctx, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// MenuCreate godoc
// @Summary 创建菜单
// @Schemes
// @Description 创建新的菜单
// @Tags 菜单模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.MenuCreateRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/menu [post]
func (h *UserHandler) MenuCreate(ctx *gin.Context) {
	var req v1.MenuCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if err := h.userService.MenuCreate(ctx, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// MenuDelete godoc
// @Summary 删除菜单
// @Schemes
// @Description 删除指定菜单
// @Tags 菜单模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query uint true "菜单ID"
// @Success 200 {object} v1.Response
// @Router /v1/admin/menu [delete]
func (h *UserHandler) MenuDelete(ctx *gin.Context) {
	var req v1.MenuDeleteRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if err := h.userService.MenuDelete(ctx, req.ID); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return

	}
	v1.HandleSuccess(ctx, nil)
}

// GetAdminMenus godoc
// @Summary 获取管理员菜单
// @Schemes
// @Description 获取管理员菜单列表
// @Tags 菜单模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetMenuResponse
// @Router /v1/admin/menus [get]
func (h *UserHandler) GetAdminMenus(ctx *gin.Context) {
	data, err := h.userService.GetAdminMenus(ctx)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	// 过滤权限菜单
	v1.HandleSuccess(ctx, data)
}

func (h *UserHandler) ListRoles(ctx *gin.Context) {
}

// RoleUpdate godoc
// @Summary 更新角色
// @Schemes
// @Description 更新角色信息
// @Tags 角色模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.RoleUpdateRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/role [put]
func (h *UserHandler) RoleUpdate(ctx *gin.Context) {
	var req v1.RoleUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if err := h.userService.RoleUpdate(ctx, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// RoleCreate godoc
// @Summary 创建角色
// @Schemes
// @Description 创建新的角色
// @Tags 角色模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.RoleCreateRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/role [post]
func (h *UserHandler) RoleCreate(ctx *gin.Context) {
	var req v1.RoleCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if err := h.userService.RoleCreate(ctx, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// RoleDelete godoc
// @Summary 删除角色
// @Schemes
// @Description 删除指定角色
// @Tags 角色模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query uint true "角色ID"
// @Success 200 {object} v1.Response
// @Router /v1/admin/role [delete]
func (h *UserHandler) RoleDelete(ctx *gin.Context) {
	var req v1.RoleDeleteRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if err := h.userService.RoleDelete(ctx, req.ID); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// GetApis godoc
// @Summary 获取API列表
// @Schemes
// @Description 获取API列表
// @Tags API模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Param group query string false "API分组"
// @Param name query string false "API名称"
// @Param path query string false "API路径"
// @Param method query string false "请求方法"
// @Success 200 {object} v1.GetApisResponse
// @Router /v1/admin/apis [get]
func (h *UserHandler) GetApis(ctx *gin.Context) {
	var req v1.GetApisRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	data, err := h.userService.ListApis(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, data)
}

// ApiUpdate godoc
// @Summary 更新API
// @Schemes
// @Description 更新API信息
// @Tags API模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.ApiUpdateRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/api [put]
func (h *UserHandler) ApiUpdate(ctx *gin.Context) {
	var req v1.ApiUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if err := h.userService.ApiUpdate(ctx, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// ApiCreate godoc
// @Summary 创建API
// @Schemes
// @Description 创建新的API
// @Tags API模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.ApiCreateRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/api [post]
func (h *UserHandler) ApiCreate(ctx *gin.Context) {
	var req v1.ApiCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if err := h.userService.ApiCreate(ctx, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// ApiDelete godoc
// @Summary 删除API
// @Schemes
// @Description 删除指定API
// @Tags API模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query uint true "API ID"
// @Success 200 {object} v1.Response
// @Router /v1/admin/api [delete]
func (h *UserHandler) ApiDelete(ctx *gin.Context) {
	var req v1.ApiDeleteRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if err := h.userService.ApiDelete(ctx, req.ID); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// GetUserPermissions godoc
// @Summary 获取用户权限
// @Schemes
// @Description 获取当前用户的权限列表
// @Tags 权限模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetUserPermissionsData
// @Router /v1/admin/user/permissions [get]
func (h *UserHandler) GetUserPermissions(ctx *gin.Context) {
	data, err := h.userService.GetUserPermissions(ctx, GetUserIdFromCtx(ctx))
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	// 过滤权限菜单
	v1.HandleSuccess(ctx, data)
}

// GetRolePermissions godoc
// @Summary 获取角色权限
// @Schemes
// @Description 获取指定角色的权限列表
// @Tags 权限模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param role query string true "角色名称"
// @Success 200 {object} v1.GetRolePermissionsData
// @Router /v1/admin/role/permissions [get]
func (h *UserHandler) GetRolePermissions(ctx *gin.Context) {
	var req v1.GetRolePermissionsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	data, err := h.userService.GetRolePermissions(ctx, req.Role)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, data)
}

// UpdateRolePermission godoc
// @Summary 更新角色权限
// @Schemes
// @Description 更新指定角色的权限列表
// @Tags 权限模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateRolePermissionRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/role/permissions [put]
func (h *UserHandler) UpdateRolePermission(ctx *gin.Context) {
	var req v1.UpdateRolePermissionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	err := h.userService.UpdateRolePermission(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
