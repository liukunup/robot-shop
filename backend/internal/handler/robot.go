package handler

import (
	"github.com/gin-gonic/gin"
	"backend/internal/service"
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

func (h *RobotHandler) GetRobot(ctx *gin.Context) {

}
