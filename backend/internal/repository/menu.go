package repository

import (
	"backend/internal/model"
	"context"
)

type MenuRepository interface {
	ListMenus(ctx context.Context) ([]model.Menu, int64, error)
	MenuCreate(ctx context.Context, m *model.Menu) error
	MenuUpdate(ctx context.Context, id uint, data map[string]interface{}) error
	MenuDelete(ctx context.Context, id uint) error
	GetMenu(ctx context.Context, id uint) (model.Menu, error)
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

func (r *menuRepository) ListMenus(ctx context.Context) ([]model.Menu, int64, error) {
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

func (r *menuRepository) MenuCreate(ctx context.Context, m *model.Menu) error {
	return r.DB(ctx).Save(m).Error
}

func (r *menuRepository) MenuUpdate(ctx context.Context, id uint, data map[string]interface{}) error {
	return r.DB(ctx).Model(&model.Menu{}).Where("id = ?", id).Updates(data).Error
}

func (r *menuRepository) MenuDelete(ctx context.Context, id uint) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&model.Menu{}).Error
}

func (r *menuRepository) GetMenu(ctx context.Context, id uint) (model.Menu, error) {
	m := model.Menu{}
	return m, r.DB(ctx).Where("id = ?", id).First(&m).Error
}
