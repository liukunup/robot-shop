package repository

import (
	v1 "backend/api/v1"
	"backend/internal/model"
	"context"
)

type RobotRepository interface {
	ListRobots(ctx context.Context, req *v1.RobotSearchRequest) ([]model.Robot, int64, error)
	RobotCreate(ctx context.Context, robot *model.Robot) error
	RobotUpdate(ctx context.Context, id uint, data map[string]interface{}) error
	RobotDelete(ctx context.Context, id uint) error
	GetRobot(ctx context.Context, id uint) (model.Robot, error)
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

func (r *robotRepository) ListRobots(ctx context.Context, req *v1.RobotSearchRequest) ([]model.Robot, int64, error) {
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

func (r *robotRepository) RobotCreate(ctx context.Context, m *model.Robot) error {
	return r.DB(ctx).Create(m).Error
}

func (r *robotRepository) RobotUpdate(ctx context.Context, id uint, data map[string]interface{}) error {
	return r.DB(ctx).Model(&model.Robot{}).Where("id = ?", id).Updates(data).Error
}

func (r *robotRepository) RobotDelete(ctx context.Context, id uint) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&model.Robot{}).Error
}

func (r *robotRepository) GetRobot(ctx context.Context, id uint) (model.Robot, error) {
	m := model.Robot{}
	return m, r.DB(ctx).Where("id = ?", id).First(&m).Error
}
