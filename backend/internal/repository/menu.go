package repository

import (
	v1 "backend/api/v1"
	"backend/internal/model"
	"context"
)

type MenuRepository interface {
	Get(ctx context.Context, id uint) (model.Menu, error)
	List(ctx context.Context, req *v1.MenuListRequest) ([]model.Menu, int64, error)
	Create(ctx context.Context, m *model.Menu) error
	Update(ctx context.Context, id uint, data map[string]interface{}) error
	Delete(ctx context.Context, id uint) error
}

func NewMenuRepository(
	repository *Repository,
) MenuRepository {
	return &menuRepository{
		Repository: repository,
	}
}

type menuRepository struct {
	*Repository
}

func (r *menuRepository) Get(ctx context.Context, id uint) (model.Menu, error) {
	m := model.Menu{}
	return m, r.DB(ctx).Where("id = ?", id).First(&m).Error
}

func (r *menuRepository) List(ctx context.Context, req *v1.MenuListRequest) ([]model.Menu, int64, error) {
	var list []model.Menu
	var total int64
	scope := r.DB(ctx).Model(&model.Menu{})
	if err := scope.Count(&total).Error; err != nil {
		return nil, total, err
	}
	if nil != req {
		scope = scope.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize)
	}
	if err := scope.Order("`weight` DESC").Find(&list).Error; err != nil {
		return nil, total, err
	}
	return list, total, nil
}

func (r *menuRepository) Create(ctx context.Context, m *model.Menu) error {
	return r.DB(ctx).Save(m).Error
}

func (r *menuRepository) Update(ctx context.Context, id uint, data map[string]interface{}) error {
	return r.DB(ctx).Model(&model.Menu{}).Where("id = ?", id).Updates(data).Error
}

func (r *menuRepository) Delete(ctx context.Context, id uint) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&model.Menu{}).Error
}
