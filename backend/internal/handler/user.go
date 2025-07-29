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
// @Description 搜索时支持邮箱、用户名、昵称字段的筛选
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "页码"
// @Param pageSize query int true "分页大小"
// @Param email query string false "邮箱"
// @Param username query string false "用户名"
// @Param nickname query string false "昵称"
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
// @Param request body v1.UserRequest true "用户数据"
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
// @Description 更新指定`ID`的用户
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path uint true "UID"
// @Param request body v1.UserRequest true "用户数据"
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
// @Description 删除指定`ID`的用户
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query uint true "UID"
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

// GetUserByID godoc
// @Summary 获取指定用户
// @Schemes
// @Description 获取指定用户的详细信息
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path uint true "UID"
// @Success 200 {object} v1.UserResponse
// @Router /users/{id} [get]
// @ID GetUserByID
func (h *UserHandler) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	uid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("GetUserByID parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	data, err := h.userService.Get(ctx, uint(uid))
	if err != nil {
		h.logger.WithContext(ctx).Error("userService.Get error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, data)
}

// GetProfile godoc
// @Summary 获取当前用户
// @Schemes
// @Description 获取当前用户的详细信息
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.UserResponse
// @Router /users/profile [get]
// @ID FetchCurrentUser
func (h *UserHandler) GetProfile(ctx *gin.Context) {
	uid := GetUserIdFromCtx(ctx)
	if uid == 0 {
		h.logger.WithContext(ctx).Error("GetProfile get uid error")
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

// UpdateProfile godoc
// @Summary 更新当前用户
// @Schemes
// @Description 更新当前用户的详细信息
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UserRequest true "用户数据"
// @Success 200 {object} v1.Response
// @Router /users/profile [put]
// @ID UpdateProfile
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	uid := GetUserIdFromCtx(ctx)
	if uid == 0 {
		h.logger.WithContext(ctx).Error("UpdateProfile get uid error")
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
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

// UploadAvatar godoc
// @Summary 上传头像
// @Schemes
// @Description 上传图片来设置或更新当前用户的头像
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param file formData file true "头像文件"
// @Success 200 {object} v1.Response
// @Router /users/profile/avatar [put]
// @ID UploadAvatar
func (h *UserHandler) UploadAvatar(ctx *gin.Context) {
	uid := GetUserIdFromCtx(ctx)
	if uid == 0 {
		h.logger.WithContext(ctx).Error("UploadAvatar get uid error")
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		h.logger.WithContext(ctx).Error("UploadAvatar get file error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	reader, err := file.Open()
	if err != nil {
		h.logger.WithContext(ctx).Error("UploadAvatar get file error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	defer reader.Close()

	req := &v1.AvatarRequest{
		UserID:   uid,
		Filename: file.Filename,
		Size:     file.Size,
		Type:     file.Header.Get("Content-Type"),
	}

	if err := h.userService.UploadAvatar(ctx, uid, req, reader); err != nil {
		h.logger.WithContext(ctx).Error("userService.UploadAvatar error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// GetMenus godoc
// @Summary 获取用户菜单
// @Schemes
// @Description 获取当前用户的菜单列表
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.DynamicMenuResponse
// @Router /users/menus [get]
// @ID FetchCurrentMenus
func (h *UserHandler) GetMenus(ctx *gin.Context) {
	uid := GetUserIdFromCtx(ctx)
	if uid == 0 {
		h.logger.WithContext(ctx).Error("GetMenu get uid error")
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	data, err := h.userService.GetMenuTree(ctx, uid)
	if err != nil {
		h.logger.WithContext(ctx).Error("userService.GetMenu error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, data)
}

// Register godoc
// @Summary 用户注册
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
// @Summary 用户登录
// @Schemes
// @Description 支持 用户名或邮箱 + 密码 进行登录
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

	tokenPair, err := h.userService.Login(ctx, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("userService.Login error", zap.Error(err))
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, tokenPair)
}

// RefreshToken godoc
// @Summary 刷新令牌
// @Schemes
// @Description 刷新访问令牌和刷新令牌，采用双Token窗口刷新机制进行滚动
// @Tags User
// @Accept json
// @Produce json
// @Param request body v1.RefreshTokenRequest true "令牌信息"
// @Success 200 {object} v1.LoginResponse
// @Router /refresh-token [post]
// @ID RefreshToken
func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	var req v1.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.WithContext(ctx).Error("RefreshToken bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	tokenPair, err := h.userService.RefreshToken(ctx, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("userService.RefreshToken error", zap.Error(err))
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, tokenPair)
}

// UpdatePassword godoc
// @Summary 更新密码
// @Schemes
// @Description 使用旧密码来更新密码
// @Tags User
// @Accept json
// @Produce json
// @Param request body v1.UpdatePasswordRequest true "密码信息"
// @Security Bearer
// @Success 200 {object} v1.Response
// @Router /users/password [put]
// @ID UpdatePassword
func (h *UserHandler) UpdatePassword(ctx *gin.Context) {
	uid := GetUserIdFromCtx(ctx)
	if uid == 0 {
		h.logger.WithContext(ctx).Error("UpdatePassword get uid error")
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	var req v1.UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.WithContext(ctx).Error("UpdatePassword bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.UpdatePassword(ctx, uid, &req); err != nil {
		h.logger.WithContext(ctx).Error("userService.UpdatePassword error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// ResetPassword godoc
// @Summary 重置密码
// @Schemes
// @Description 通过邮箱来进行密码重置
// @Tags User
// @Accept json
// @Produce json
// @Param request body v1.ResetPasswordRequest true "重置请求的验证信息"
// @Success 200 {object} v1.Response
// @Router /reset-password [post]
// @ID ResetPassword
func (h *UserHandler) ResetPassword(ctx *gin.Context) {
	var req v1.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.WithContext(ctx).Error("ResetPassword bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.ResetPassword(ctx, &req); err != nil {
		h.logger.WithContext(ctx).Error("userService.ResetPassword error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}
