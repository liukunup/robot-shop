package repository

import (
	v1 "backend/api/v1"
	"backend/internal/model"
	"context"
)

type RobotRepository interface {
	Get(ctx context.Context, id uint) (model.Robot, error)
	List(ctx context.Context, req *v1.RobotSearchRequest) ([]model.Robot, int64, error)
	Create(ctx context.Context, v *model.Robot) error
	Update(ctx context.Context, id uint, data map[string]interface{}) error
	Delete(ctx context.Context, id uint) error
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

func (r *robotRepository) Get(ctx context.Context, id uint) (model.Robot, error) {
	m := model.Robot{}
	return m, r.DB(ctx).Where("id = ?", id).First(&m).Error
}

func (r *robotRepository) List(ctx context.Context, req *v1.RobotSearchRequest) ([]model.Robot, int64, error) {
	var list []model.Robot
	var total int64
	scope := r.DB(ctx).Model(&model.Robot{})
	if req.Name != "" {
		scope = scope.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Desc != "" {
		scope = scope.Where("desc LIKE ?", "%"+req.Desc+"%")
	}
	if req.Owner != "" {
		scope = scope.Where("owner = ?", req.Owner)
	}
	if err := scope.Count(&total).Error; err != nil {
		return nil, total, err
	}
	if err := scope.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error; err != nil {
		return nil, total, err
	}
	return list, total, nil
}

func (r *robotRepository) Create(ctx context.Context, v *model.Robot) error {
	return r.DB(ctx).Create(v).Error
}

func (r *robotRepository) Update(ctx context.Context, id uint, data map[string]interface{}) error {
	return r.DB(ctx).Model(&model.Robot{}).Where("id = ?", id).Updates(data).Error
}

func (r *robotRepository) Delete(ctx context.Context, id uint) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&model.Robot{}).Error
}
