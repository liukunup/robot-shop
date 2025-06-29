package service

import (
	v1 "backend/api/v1"
	"backend/internal/constant"
	"backend/internal/model"
	"backend/internal/repository"
	"context"
)

type ApiService interface {
	ListApis(ctx context.Context, req *v1.ApiSearchRequest) (*v1.ApiSearchResponseData, error)
	ApiCreate(ctx context.Context, req *v1.ApiRequest) error
	ApiUpdate(ctx context.Context, id uint, req *v1.ApiRequest) error
	ApiDelete(ctx context.Context, id uint) error
	GetApi(ctx context.Context, id uint) (model.Api, error)
}

func NewApiService(
	service *Service,
	apiRepository repository.ApiRepository,
) ApiService {
	return &apiService{
		Service:       service,
		apiRepository: apiRepository,
	}
}

type apiService struct {
	*Service
	apiRepository repository.ApiRepository
}

func (s *apiService) ListApis(ctx context.Context, req *v1.ApiSearchRequest) (*v1.ApiSearchResponseData, error) {
	list, total, err := s.apiRepository.ListApis(ctx, req)
	if err != nil {
		return nil, err
	}
	groups, err := s.apiRepository.GetGroups(ctx)
	if err != nil {
		return nil, err
	}
	data := &v1.ApiSearchResponseData{
		List:   make([]v1.ApiDataItem, 0),
		Total:  total,
		Groups: groups,
	}
	for _, api := range list {
		data.List = append(data.List, v1.ApiDataItem{
			ID:        api.ID,
			CreatedAt: api.CreatedAt.Format(constant.DateTimeLayout),
			UpdatedAt: api.UpdatedAt.Format(constant.DateTimeLayout),
			Group:     api.Group,
			Method:    api.Method,
			Name:      api.Name,
			Path:      api.Path,
		})
	}
	return data, nil
}

func (s *apiService) ApiCreate(ctx context.Context, req *v1.ApiRequest) error {
	return s.apiRepository.ApiCreate(ctx, &model.Api{
		Group:  req.Group,
		Name:   req.Name,
		Path:   req.Path,
		Method: req.Method,
	})
}

func (s *apiService) ApiUpdate(ctx context.Context, id uint, req *v1.ApiRequest) error {
	data := map[string]interface{}{
		"group":  req.Group,
		"name":   req.Name,
		"path":   req.Path,
		"method": req.Method,
	}
	return s.apiRepository.ApiUpdate(ctx, id, data)
}

func (s *apiService) ApiDelete(ctx context.Context, id uint) error {
	return s.apiRepository.ApiDelete(ctx, id)
}

func (s *apiService) GetApi(ctx context.Context, id uint) (model.Api, error) {
	return s.apiRepository.GetApi(ctx, id)
}
