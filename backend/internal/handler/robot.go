package handler

import (
	v1 "backend/api/v1"
	"backend/internal/service"
	"backend/pkg/time"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

// GetRobot godoc
// @Summary 获取机器人
// @Schemes
// @Description
// @Tags 机器人模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "机器人ID"
// @Success 200 {object} v1.RobotResponseData
// @Router /robot/{id} [get]
func (h *RobotHandler) GetRobot(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	robot, err := h.robotService.GetRobot(ctx, id)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}

	v1.HandleSuccess(ctx, v1.RobotResponseData{
		Id:        robot.Id,
		RobotId:   robot.RobotId,
		Name:      robot.Name,
		Desc:      robot.Desc,
		Webhook:   robot.Webhook,
		Callback:  robot.Callback,
		Options:   robot.Options,
		Enabled:   robot.Enabled,
		Owner:     robot.Owner,
		CreatedAt: time.FormatTime(robot.CreatedAt),
		UpdatedAt: time.FormatTime(robot.UpdatedAt),
	})
}

// CreateRobot godoc
// @Summary 创建机器人
// @Schemes
// @Description
// @Tags 机器人模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.RobotRequest true "params"
// @Success 200 {object} v1.RobotResponseData
// @Router /robot [post]
func (h *RobotHandler) CreateRobot(ctx *gin.Context) {
	var req v1.RobotRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid body"})
		return
	}

	robot, err := h.robotService.CreateRobot(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}

	v1.HandleSuccess(ctx, v1.RobotResponseData{
		Id:        robot.Id,
		RobotId:   robot.RobotId,
		Name:      robot.Name,
		Desc:      robot.Desc,
		Webhook:   robot.Webhook,
		Callback:  robot.Callback,
		Options:   robot.Options,
		Enabled:   robot.Enabled,
		Owner:     robot.Owner,
		CreatedAt: time.FormatTime(robot.CreatedAt),
		UpdatedAt: time.FormatTime(robot.UpdatedAt),
	})
}

// UpdateRobot godoc
// @Summary 更新机器人
// @Schemes
// @Description
// @Tags 机器人模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "机器人ID"
// @Param request body v1.RobotRequest true "params"
// @Success 200 {object} v1.RobotResponseData
// @Router /robot/{id} [put]
func (h *RobotHandler) UpdateRobot(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req v1.RobotRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid body"})
		return
	}

	robot, err := h.robotService.UpdateRobot(ctx, id, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}

	v1.HandleSuccess(ctx, v1.RobotResponseData{
		Id:        robot.Id,
		RobotId:   robot.RobotId,
		Name:      robot.Name,
		Desc:      robot.Desc,
		Webhook:   robot.Webhook,
		Callback:  robot.Callback,
		Options:   robot.Options,
		Enabled:   robot.Enabled,
		Owner:     robot.Owner,
		CreatedAt: time.FormatTime(robot.CreatedAt),
		UpdatedAt: time.FormatTime(robot.UpdatedAt),
	})
}

// DeleteRobot godoc
// @Summary 删除机器人
// @Schemes
// @Description
// @Tags 机器人模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "机器人ID"
// @Success 200 {object} v1.Response
// @Router /robot/{id} [delete]
func (h *RobotHandler) DeleteRobot(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.robotService.DeleteRobot(ctx, id); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// ListRobots godoc
// @Summary 获取机器人列表
// @Schemes
// @Description
// @Tags 机器人模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {object} v1.PageResponse[v1.RobotResponseData]
// @Router /robot [get]
func (h *RobotHandler) ListRobots(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid page param"})
		return
	}

	size, err := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, gin.H{"error": "invalid size param"})
		return
	}

	options := make(map[string]interface{})
	for key, value := range ctx.Request.URL.Query() {
		if key == "page" || key == "size" {
			continue
		}
		options[key] = value[0]
	}

	robots, total, err := h.robotService.ListRobots(ctx, page, size, options)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]v1.RobotResponseData, 0)
	for _, robot := range robots {
		items = append(items, v1.RobotResponseData{
			Id:        robot.Id,
			RobotId:   robot.RobotId,
			Name:      robot.Name,
			Desc:      robot.Desc,
			Webhook:   robot.Webhook,
			Callback:  robot.Callback,
			Options:   robot.Options,
			Enabled:   robot.Enabled,
			Owner:     robot.Owner,
			CreatedAt: time.FormatTime(robot.CreatedAt),
			UpdatedAt: time.FormatTime(robot.UpdatedAt),
		})
	}

	v1.HandleSuccess(ctx, v1.PageResponse[v1.RobotResponseData]{
		Items: items,
		Total: total,
	})
}
