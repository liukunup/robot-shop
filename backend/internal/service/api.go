package service

import (
	v1 "backend/api/v1"
	"backend/internal/constant"
	"backend/internal/model"
	"backend/internal/repository"
	"context"
)

type ApiService interface {
	Get(ctx context.Context, id uint) (model.Api, error)
	List(ctx context.Context, req *v1.ApiSearchRequest) (*v1.ApiSearchResponseData, error)
	Create(ctx context.Context, req *v1.ApiRequest) error
	Update(ctx context.Context, id uint, req *v1.ApiRequest) error
	Delete(ctx context.Context, id uint) error
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

func (s *apiService) Get(ctx context.Context, id uint) (model.Api, error) {
	return s.apiRepository.Get(ctx, id)
}

func (s *apiService) List(ctx context.Context, req *v1.ApiSearchRequest) (*v1.ApiSearchResponseData, error) {
	list, total, err := s.apiRepository.List(ctx, req)
	if err != nil {
		return nil, err
	}
	data := &v1.ApiSearchResponseData{
		List:   make([]v1.ApiDataItem, 0),
		Total:  total,
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

func (s *apiService) Create(ctx context.Context, req *v1.ApiRequest) error {
	return s.apiRepository.Create(ctx, &model.Api{
		Group:  req.Group,
		Name:   req.Name,
		Path:   req.Path,
		Method: req.Method,
	})
}

func (s *apiService) Update(ctx context.Context, id uint, req *v1.ApiRequest) error {
	data := map[string]interface{}{
		"group":  req.Group,
		"name":   req.Name,
		"path":   req.Path,
		"method": req.Method,
	}
	return s.apiRepository.Update(ctx, id, data)
}

func (s *apiService) Delete(ctx context.Context, id uint) error {
	return s.apiRepository.Delete(ctx, id)
}
