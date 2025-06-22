package repository

import (
	v1 "backend/api/v1"
	"backend/internal/model"
	"context"
)

type RobotRepository interface {
	List(ctx context.Context, req *v1.ListRobotRequest) ([]model.Robot, int64, error)
	Create(ctx context.Context, robot *model.Robot) error
	Update(ctx context.Context, robot *model.Robot) error
	Delete(ctx context.Context, id uint) error
	Get(ctx context.Context, id uint) (model.Robot, error)
}

func NewRobotRepository(
	repository *Repository,
) RobotRepository {
	return &robotRepository{
		Repository: repository,
	}
}

type robotRepository struct {
	*Repository
}

func (r *robotRepository) List(ctx context.Context, req *v1.ListRobotRequest) ([]model.Robot, int64, error) {
	var list []model.Robot
	var total int64
	scope := r.DB(ctx).Model(&model.Robot{})
	if req.Name != "" {
		scope = scope.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if err := scope.Count(&total).Error; err != nil {
		return nil, total, err
	}
	if err := scope.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error; err != nil {
		return nil, total, err
	}
	return list, total, nil
}

func (r *robotRepository) Create(ctx context.Context, m *model.Robot) error {
	return r.DB(ctx).Create(m).Error
}

func (r *robotRepository) Update(ctx context.Context, m *model.Robot) error {
	return r.DB(ctx).Where("id = ?", m.ID).Save(m).Error
}

func (r *robotRepository) Delete(ctx context.Context, id uint) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&model.Robot{}).Error
}

func (r *robotRepository) Get(ctx context.Context, uid uint) (model.Robot, error) {
	m := model.Robot{}
	return m, r.DB(ctx).Where("id = ?", uid).First(&m).Error
}
