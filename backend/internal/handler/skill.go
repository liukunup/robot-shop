package handler

import (
	"github.com/gin-gonic/gin"
	"backend/internal/service"
)

type SkillHandler struct {
	*Handler
	skillService service.SkillService
}

func NewSkillHandler(
    handler *Handler,
    skillService service.SkillService,
) *SkillHandler {
	return &SkillHandler{
		Handler:      handler,
		skillService: skillService,
	}
}

func (h *SkillHandler) GetSkill(ctx *gin.Context) {

}
