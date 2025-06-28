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

type RobotHandler struct {
	*Handler
	robotService service.RobotService
}

func NewRobotHandler(
	handler *Handler,
	robotService service.RobotService,
) *RobotHandler {
	return &RobotHandler{
		Handler:      handler,
		robotService: robotService,
	}
}

// ListRobots godoc
// @Summary 获取机器人列表
// @Schemes
// @Description 搜索时支持名称、描述和所有者筛选
// @Tags Robot
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "页码"
// @Param pageSize query int true "分页大小"
// @Param name query string false "名称"
// @Param desc query string false "描述"
// @Param owner query string false "所有者"
// @Success 200 {object} v1.RobotSearchResponse
// @Router /robots [get]
// @ID ListRobots
func (h *RobotHandler) ListRobots(ctx *gin.Context) {
	var req v1.RobotSearchRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.WithContext(ctx).Error("ListRobots bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	data, err := h.robotService.List(ctx, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("robotService.List error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, data)
}

// RobotCreate godoc
// @Summary 创建机器人
// @Schemes
// @Description 创建一个新的机器人
// @Tags Robot
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.RobotRequest true "机器人数据"
// @Success 200 {object} v1.Response
// @Router /robots [post]
// @ID RobotCreate
func (h *RobotHandler) RobotCreate(ctx *gin.Context) {
	var req v1.RobotRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.WithContext(ctx).Error("RobotCreate bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	err := h.robotService.Create(ctx, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("robotService.Create error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// RobotUpdate godoc
// @Summary 更新机器人
// @Schemes
// @Description 更新机器人配置
// @Tags Robot
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path uint true "机器人ID"
// @Param request body v1.RobotRequest true "机器人数据"
// @Success 200 {object} v1.Response
// @Router /robots/{id} [put]
// @ID RobotUpdate
func (h *RobotHandler) RobotUpdate(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("RobotUpdate parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req v1.RobotRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.WithContext(ctx).Error("RobotUpdate bind error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.robotService.Update(ctx, uint(id), &req); err != nil {
		h.logger.WithContext(ctx).Error("robotService.Update error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// RobotDelete godoc
// @Summary 删除机器人
// @Schemes
// @Description
// @Tags Robot
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path uint true "机器人ID"
// @Success 200 {object} v1.Response
// @Router /robots/{id} [delete]
// @ID RobotDelete
func (h *RobotHandler) RobotDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("RobotDelete parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.robotService.Delete(ctx, uint(id)); err != nil {
		h.logger.WithContext(ctx).Error("robotService.Delete error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// GetRobot godoc
// @Summary 获取机器人详情
// @Schemes
// @Description
// @Tags Robot
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path uint true "机器人ID"
// @Success 200 {object} v1.RobotResponse
// @Router /robots/{id} [get]
// @ID GetRobot
func (h *RobotHandler) GetRobot(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.WithContext(ctx).Error("GetRobot parse id error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	robot, err := h.robotService.Get(ctx, uint(id))
	if err != nil {
		h.logger.WithContext(ctx).Error("robotService.Get error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}

	v1.HandleSuccess(ctx, v1.RobotDataItem{
		Id:        robot.ID,
		CreatedAt: time.FormatTime(robot.CreatedAt),
		UpdatedAt: time.FormatTime(robot.UpdatedAt),
		Name:      robot.Name,
		Desc:      robot.Desc,
		Webhook:   robot.Webhook,
		Callback:  robot.Callback,
		Enabled:   robot.Enabled,
		Owner:     robot.Owner,
	})
}
