package repository

import (
    "context"
	"backend/internal/model"
)

type SkillRepository interface {
	GetSkill(ctx context.Context, id int64) (*model.Skill, error)
}

func NewSkillRepository(
	repository *Repository,
) SkillRepository {
	return &skillRepository{
		Repository: repository,
	}
}

type skillRepository struct {
	*Repository
}

func (r *skillRepository) GetSkill(ctx context.Context, id int64) (*model.Skill, error) {
	var skill model.Skill

	return &skill, nil
}
