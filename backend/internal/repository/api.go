package repository

import (
	v1 "backend/api/v1"
	"backend/internal/model"
	"context"
)

type ApiRepository interface {
	GetApis(ctx context.Context, req *v1.GetApisRequest) ([]model.Api, int64, error)
	GetApiGroups(ctx context.Context) ([]string, error)
	ApiUpdate(ctx context.Context, m *model.Api) error
	ApiCreate(ctx context.Context, m *model.Api) error
	ApiDelete(ctx context.Context, id uint) error
}

func NewApiRepository(
	repository *Repository,
) ApiRepository {
	return &apiRepository{
		Repository: repository,
	}
}

type apiRepository struct {
	*Repository
}

func (r *apiRepository) GetApis(ctx context.Context, req *v1.GetApisRequest) ([]model.Api, int64, error) {
	var list []model.Api
	var total int64
	scope := r.DB(ctx).Model(&model.Api{})
	if req.Name != "" {
		scope = scope.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Group != "" {
		scope = scope.Where("`group` LIKE ?", "%"+req.Group+"%")
	}
	if req.Path != "" {
		scope = scope.Where("path LIKE ?", "%"+req.Path+"%")
	}
	if req.Method != "" {
		scope = scope.Where("method = ?", req.Method)
	}
	if err := scope.Count(&total).Error; err != nil {
		return nil, total, err
	}
	if err := scope.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("`group` ASC").Find(&list).Error; err != nil {
		return nil, total, err
	}
	return list, total, nil
}

func (r *apiRepository) GetApiGroups(ctx context.Context) ([]string, error) {
	res := make([]string, 0)
	if err := r.DB(ctx).Model(&model.Api{}).Group("`group`").Pluck("`group`", &res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *apiRepository) ApiUpdate(ctx context.Context, m *model.Api) error {
	return r.DB(ctx).Where("id = ?", m.ID).Save(m).Error
}

func (r *apiRepository) ApiCreate(ctx context.Context, m *model.Api) error {
	return r.DB(ctx).Create(m).Error
}

func (r *apiRepository) ApiDelete(ctx context.Context, id uint) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&model.Api{}).Error
}
