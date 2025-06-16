package repository

import (
	"backend/internal/model"
	"context"
)

type MenuRepository interface {
	ListMenus(ctx context.Context) ([]model.Menu, error)
	MenuUpdate(ctx context.Context, m *model.Menu) error
	MenuCreate(ctx context.Context, m *model.Menu) error
	MenuDelete(ctx context.Context, id uint) error
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

func (r *menuRepository) ListMenus(ctx context.Context) ([]model.Menu, error) {
	var menuList []model.Menu
	if err := r.DB(ctx).Order("weight DESC").Find(&menuList).Error; err != nil {
		return nil, err
	}
	return menuList, nil
}

func (r *menuRepository) MenuUpdate(ctx context.Context, m *model.Menu) error {
	return r.DB(ctx).Where("id = ?", m.ID).Save(m).Error
}

func (r *menuRepository) MenuCreate(ctx context.Context, m *model.Menu) error {
	return r.DB(ctx).Save(m).Error
}

func (r *menuRepository) MenuDelete(ctx context.Context, id uint) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&model.Menu{}).Error
}
