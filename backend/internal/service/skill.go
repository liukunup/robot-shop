package service

import (
    "context"
	"backend/internal/model"
	"backend/internal/repository"
)

type SkillService interface {
	GetSkill(ctx context.Context, id int64) (*model.Skill, error)
}
func NewSkillService(
    service *Service,
    skillRepository repository.SkillRepository,
) SkillService {
	return &skillService{
		Service:        service,
		skillRepository: skillRepository,
	}
}

type skillService struct {
	*Service
	skillRepository repository.SkillRepository
}

func (s *skillService) GetSkill(ctx context.Context, id int64) (*model.Skill, error) {
	return s.skillRepository.GetSkill(ctx, id)
}
