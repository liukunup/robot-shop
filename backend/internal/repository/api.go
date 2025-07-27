package repository

import (
	v1 "backend/api/v1"
	"backend/internal/model"
	"context"
)

type ApiRepository interface {
	Get(ctx context.Context, id uint) (model.Api, error)
	List(ctx context.Context, req *v1.ApiSearchRequest) ([]model.Api, int64, error)
	Create(ctx context.Context, m *model.Api) error
	Update(ctx context.Context, id uint, data map[string]interface{}) error
	Delete(ctx context.Context, id uint) error

	ListAllGroups(ctx context.Context) ([]string, error)
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

func (r *apiRepository) Get(ctx context.Context, id uint) (model.Api, error) {
	m := model.Api{}
	return m, r.DB(ctx).Where("id = ?", id).First(&m).Error
}

func (r *apiRepository) List(ctx context.Context, req *v1.ApiSearchRequest) ([]model.Api, int64, error) {
	var list []model.Api
	var total int64
	scope := r.DB(ctx).Model(&model.Api{})
	if req.Group != "" {
		scope = scope.Where("`group` LIKE ?", "%"+req.Group+"%")
	}
	if req.Name != "" {
		scope = scope.Where("name LIKE ?", "%"+req.Name+"%")
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

func (r *apiRepository) Create(ctx context.Context, m *model.Api) error {
	return r.DB(ctx).Create(m).Error
}

func (r *apiRepository) Update(ctx context.Context, id uint, data map[string]interface{}) error {
	return r.DB(ctx).Model(&model.Api{}).Where("id = ?", id).Updates(data).Error
}

func (r *apiRepository) Delete(ctx context.Context, id uint) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&model.Api{}).Error
}

func (r *apiRepository) ListAllGroups(ctx context.Context) ([]string, error) {
	groups := make([]string, 0)
	if err := r.DB(ctx).Model(&model.Api{}).Group("`group`").Pluck("`group`", &groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}
