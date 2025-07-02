package handler

import (
	v1 "backend/api/v1"
	"backend/internal/service"
	"net/http"
	"strconv"

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

// ListUsers godoc
// @Summary 获取用户列表
// @Schemes
// @Description 搜索时支持用户名、昵称、手机和邮箱筛选
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "页码"
// @Param pageSize query int true "分页大小"
// @Param username query string false "用户名"
// @Param nickname query string false "昵称"
// @Param phone query string false "手机"
// @Param email query string false "邮箱"
// @Success 200 {object} v1.UserSearchResponse
// @Router /admin/users [get]
// @ID ListUsers
func (h *UserHandler) ListUsers(ctx *gin.Context) {
	var req v1.UserSearchRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("ListUsers bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	data, err := h.userService.List(ctx, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("userService.List error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, data)
}

// CreateUser godoc
// @Summary 创建用户
// @Schemes
// @Description 创建一个新的用户
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UserRequest true "用户信息"
// @Success 200 {object} v1.Response
// @Router /admin/users [post]
// @ID CreateUser
func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var req v1.UserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("CreateUser bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.Create(ctx, &req); err != nil {
		h.logger.WithContext(ctx).Error("userService.Create error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// UpdateUser godoc
// @Summary 更新用户
// @Schemes
// @Description 更新用户信息
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path uint true "用户ID"
// @Param request body v1.UserRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /admin/users/{id} [put]
// @ID UpdateUser
func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	uid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("UpdateUser parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req v1.UserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("UpdateUser bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.Update(ctx, uint(uid), &req); err != nil {
		h.logger.WithContext(ctx).Error("userService.Update error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// DeleteUser godoc
// @Summary 删除用户
// @Schemes
// @Description 删除指定ID的用户
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query uint true "用户ID"
// @Success 200 {object} v1.Response
// @Router /admin/users/{id} [delete]
// @ID DeleteUser
func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	uid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("DeleteUser parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.userService.Delete(ctx, uint(uid)); err != nil {
		h.logger.WithContext(ctx).Error("userService.Delete error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// GetCurrentUser godoc
// @Summary 获取当前用户
// @Schemes
// @Description 获取当前用户的详细信息
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.UserResponse
// @Router /users/me [get]
// @ID GetCurrentUser
func (h *UserHandler) GetCurrentUser(ctx *gin.Context) {
	uid := GetUserIdFromCtx(ctx)
	if uid == 0 {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	data, err := h.userService.Get(ctx, uid)
	if err != nil {
		h.logger.WithContext(ctx).Error("userService.Get error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, data)
}

// GetUserPermissions godoc
// @Summary 获取用户权限
// @Schemes
// @Description 获取当前用户的权限列表
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.UserPermissionResponse
// @Router /users/me/permissions [get]
// @ID GetUserPermissions
func (h *UserHandler) GetUserPermissions(ctx *gin.Context) {
	uid := GetUserIdFromCtx(ctx)
	if uid == 0 {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	data, err := h.userService.GetPermissions(ctx, uid)
	if err != nil {
		h.logger.WithContext(ctx).Error("userService.GetPermissions error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, data)
}

// GetUserMenus godoc
// @Summary 获取用户菜单
// @Schemes
// @Description 获取当前用户的菜单列表
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.MenuSearchResponse
// @Router /users/me/menus [get]
// @ID GetUserMenus
func (h *UserHandler) GetUserMenus(ctx *gin.Context) {
	uid := GetUserIdFromCtx(ctx)
	if uid == 0 {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	data, err := h.userService.GetMenus(ctx, uid)
	if err != nil {
		h.logger.WithContext(ctx).Error("menuService.GetMenus error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, data)
}

// Register godoc
// @Summary 注册
// @Schemes
// @Description 目前只支持通过邮箱进行注册
// @Tags User
// @Accept json
// @Produce json
// @Param request body v1.RegisterRequest true "注册信息"
// @Success 200 {object} v1.Response
// @Router /register [post]
// @ID Register
func (h *UserHandler) Register(ctx *gin.Context) {
	var req v1.RegisterRequest
	if err := ctx.ShouldBindJSON(req); err != nil {
		h.logger.WithContext(ctx).Error("Register bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.Register(ctx, &req); err != nil {
		h.logger.WithContext(ctx).Error("userService.Register error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// Login godoc
// @Summary 登录
// @Schemes
// @Description 支持用户名或邮箱登录
// @Tags User
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "登录凭证"
// @Success 200 {object} v1.LoginResponse
// @Router /login [post]
// @ID Login
func (h *UserHandler) Login(ctx *gin.Context) {
	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.WithContext(ctx).Error("Login bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	token, err := h.userService.Login(ctx, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("userService.Login error", zap.Error(err))
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, gin.H{"error": err.Error()})
		return
	}

	v1.HandleSuccess(ctx, v1.LoginResponseData{
		AccessToken: token,
	})
}
