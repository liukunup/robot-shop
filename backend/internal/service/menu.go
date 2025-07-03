package service

import (
	v1 "backend/api/v1"
	"backend/internal/constant"
	"backend/internal/model"
	"backend/internal/repository"
	"context"

	"go.uber.org/zap"
)

type MenuService interface {
	List(ctx context.Context, req *v1.MenuListRequest) (*v1.MenuListResponseData, error)
	Create(ctx context.Context, req *v1.MenuRequest) error
	Update(ctx context.Context, id uint, req *v1.MenuRequest) error
	Delete(ctx context.Context, id uint) error
}

func NewMenuService(
	service *Service,
	menuRepository repository.MenuRepository,
) MenuService {
	return &menuService{
		Service:        service,
		menuRepository: menuRepository,
	}
}

type menuService struct {
	*Service
	menuRepository repository.MenuRepository
}

func (s *menuService) List(ctx context.Context, req *v1.MenuListRequest) (*v1.MenuListResponseData, error) {
	list, total, err := s.menuRepository.List(ctx, req)
	if err != nil {
		s.logger.WithContext(ctx).Error("List error", zap.Error(err))
		return nil, err
	}
	data := &v1.MenuListResponseData{
		List:  make([]v1.MenuDataItem, 0),
		Total: total,
	}
	for _, menu := range list {
		data.List = append(data.List, v1.MenuDataItem{
			ID:        menu.ID,
			CreatedAt: menu.CreatedAt.Format(constant.DateTimeLayout),
			UpdatedAt: menu.UpdatedAt.Format(constant.DateTimeLayout),
			ParentID:  menu.ParentID,
			Path:      menu.Path,
			Redirect:  menu.Redirect,
			Component: menu.Component,
			Name:      menu.Name,
			Icon:      menu.Icon,
			Access:    menu.Access,
			Weight:    menu.Weight,
		})
	}
	return data, nil
}

func (s *menuService) Create(ctx context.Context, req *v1.MenuRequest) error {
	return s.menuRepository.Create(ctx, &model.Menu{
		ParentID:  req.ParentID,
		Path:      req.Path,
		Redirect:  req.Redirect,
		Component: req.Component,
		Name:      req.Name,
		Icon:      req.Icon,
		Access:    req.Access,
		Weight:    req.Weight,
	})
}

func (s *menuService) Update(ctx context.Context, id uint, req *v1.MenuRequest) error {
	data := map[string]interface{}{
		"parent_id": req.ParentID,
		"path":      req.Path,
		"redirect":  req.Redirect,
		"component": req.Component,
		"name":      req.Name,
		"icon":      req.Icon,
		"access":    req.Access,
		"weight":    req.Weight,
	}
	return s.menuRepository.Update(ctx, id, data)
}

func (s *menuService) Delete(ctx context.Context, id uint) error {
	return s.menuRepository.Delete(ctx, id)
}
