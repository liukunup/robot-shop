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
	ListMenus(ctx context.Context, req *v1.MenuSearchRequest) (*v1.MenuSearchResponseData, error)
	MenuCreate(ctx context.Context, req *v1.MenuRequest) error
	MenuUpdate(ctx context.Context, id uint, req *v1.MenuRequest) error
	MenuDelete(ctx context.Context, id uint) error
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

func (s *menuService) ListMenus(ctx context.Context, req *v1.MenuSearchRequest) (*v1.MenuSearchResponseData, error) {
	list, total, err := s.menuRepository.ListMenus(ctx, req)
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
			ID:         menu.ID,
			Name:       menu.Name,
			Title:      menu.Title,
			Path:       menu.Path,
			Component:  menu.Component,
			Redirect:   menu.Redirect,
			KeepAlive:  menu.KeepAlive,
			HideInMenu: menu.HideInMenu,
			Locale:     menu.Locale,
			Weight:     menu.Weight,
			Icon:       menu.Icon,
			ParentID:   menu.ParentID,
			UpdatedAt:  menu.UpdatedAt.Format(constant.DateTimeLayout),
			URL:        menu.URL,
		})
	}
	return data, nil
}

func (s *menuService) MenuCreate(ctx context.Context, req *v1.MenuRequest) error {
	return s.menuRepository.MenuCreate(ctx, &model.Menu{
		Component:  req.Component,
		Icon:       req.Icon,
		KeepAlive:  req.KeepAlive,
		HideInMenu: req.HideInMenu,
		Locale:     req.Locale,
		Weight:     req.Weight,
		Name:       req.Name,
		ParentID:   req.ParentID,
		Path:       req.Path,
		Redirect:   req.Redirect,
		Title:      req.Title,
		URL:        req.URL,
	})
}

func (s *menuService) MenuUpdate(ctx context.Context, id uint, req *v1.MenuRequest) error {
	data := map[string]interface{}{
		"component":    req.Component,
		"icon":         req.Icon,
		"keep_alive":   req.KeepAlive,
		"hide_in_menu": req.HideInMenu,
		"locale":       req.Locale,
		"weight":       req.Weight,
		"name":         req.Name,
		"parent_id":    req.ParentID,
		"path":         req.Path,
		"redirect":     req.Redirect,
		"title":        req.Title,
		"url":          req.URL,
	}
	return s.menuRepository.MenuUpdate(ctx, id, data)
}

func (s *menuService) MenuDelete(ctx context.Context, id uint) error {
	return s.menuRepository.MenuDelete(ctx, id)
}
