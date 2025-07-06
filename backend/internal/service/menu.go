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
	List(ctx context.Context, req *v1.MenuSearchRequest) (*v1.MenuSearchResponseData, error)
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

func (s *menuService) List(ctx context.Context, req *v1.MenuSearchRequest) (*v1.MenuSearchResponseData, error) {
	list, total, err := s.menuRepository.List(ctx, req)
	if err != nil {
		s.logger.WithContext(ctx).Error("List error", zap.Error(err))
		return nil, err
	}
	data := &v1.MenuSearchResponseData{
		List:  make([]v1.MenuDataItem, 0),
		Total: total,
	}
	for _, menu := range list {
		data.List = append(data.List, v1.MenuDataItem{
			ID:                 menu.ID,
			CreatedAt:          menu.CreatedAt.Format(constant.DateTimeLayout),
			UpdatedAt:          menu.UpdatedAt.Format(constant.DateTimeLayout),
			ParentID:           menu.ParentID,
			Icon:               menu.Icon,
			Name:               menu.Name,
			Path:               menu.Path,
			Component:          menu.Component,
			Access:             menu.Access,
			Locale:             menu.Locale,
			Redirect:           menu.Redirect,
			Target:             menu.Target,
			HideChildrenInMenu: menu.HideChildrenInMenu,
			HideInMenu:         menu.HideInMenu,
			FlatMenu:           menu.FlatMenu,
			Disabled:           menu.Disabled,
			Tooltip:            menu.Tooltip,
			DisabledTooltip:    menu.DisabledTooltip,
			Key:                menu.Key,
			ParentKeys:         menu.ParentKeys,
		})
	}
	return data, nil
}

func (s *menuService) Create(ctx context.Context, req *v1.MenuRequest) error {
	return s.menuRepository.Create(ctx, &model.Menu{
		ParentID:           req.ParentID,
		Icon:               req.Icon,
		Name:               req.Name,
		Path:               req.Path,
		Component:          req.Component,
		Access:             req.Access,
		Redirect:           req.Redirect,
		Target:             req.Target,
		HideChildrenInMenu: req.HideChildrenInMenu,
		HideInMenu:         req.HideInMenu,
		FlatMenu:           req.FlatMenu,
		Disabled:           req.Disabled,
		Tooltip:            req.Tooltip,
		DisabledTooltip:    req.DisabledTooltip,
		Key:                req.Key,
		ParentKeys:         req.ParentKeys,
	})
}

func (s *menuService) Update(ctx context.Context, id uint, req *v1.MenuRequest) error {
	data := map[string]interface{}{
		"parent_id":             req.ParentID,
		"icon":                  req.Icon,
		"name":                  req.Name,
		"path":                  req.Path,
		"component":             req.Component,
		"access":                req.Access,
		"redirect":              req.Redirect,
		"target":                req.Target,
		"hide_children_in_menu": req.HideChildrenInMenu,
		"hide_in_menu":          req.HideInMenu,
		"flat_menu":             req.FlatMenu,
		"disabled":              req.Disabled,
		"tooltip":               req.Tooltip,
		"disabled_tooltip":      req.DisabledTooltip,
		"key":                   req.Key,
		"parent_keys":           req.ParentKeys,
	}
	return s.menuRepository.Update(ctx, id, data)
}

func (s *menuService) Delete(ctx context.Context, id uint) error {
	return s.menuRepository.Delete(ctx, id)
}
