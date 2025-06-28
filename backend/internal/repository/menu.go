package repository

import (
	"backend/internal/model"
	"context"
)

type MenuRepository interface {
	List(ctx context.Context) ([]model.Menu, int64, error)
	Create(ctx context.Context, m *model.Menu) error
	Update(ctx context.Context, id uint, data map[string]interface{}) error
	Delete(ctx context.Context, id uint) error
	Get(ctx context.Context, id uint) (model.Menu, error)
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

func (r *menuRepository) List(ctx context.Context) ([]model.Menu, int64, error) {
	var list []model.Menu
	var total int64
	if err := r.DB(ctx).Model(&model.Menu{}).Count(&total).Error; err != nil {
		return nil, total, err
	}
	if err := r.DB(ctx).Model(&model.Menu{}).Order("`weight` DESC").Find(&list).Error; err != nil {
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

func (r *menuRepository) Get(ctx context.Context, id uint) (model.Menu, error) {
	m := model.Menu{}
	return m, r.DB(ctx).Where("id = ?", id).First(&m).Error
}
