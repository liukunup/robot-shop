package service

import (
	v1 "backend/api/v1"
	"backend/internal/constant"
	"backend/internal/model"
	"backend/internal/repository"
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type RoleService interface {
	// CRUD
	ListRoles(ctx context.Context, req *v1.RoleSearchRequest) (*v1.RoleSearchResponseData, error)
	RoleCreate(ctx context.Context, req *v1.RoleRequest) error
	RoleUpdate(ctx context.Context, id uint, req *v1.RoleRequest) error
	RoleDelete(ctx context.Context, id uint) error
	// Permission
	GetRolePermission(ctx context.Context, role string) (*v1.GetRolePermissionResponseData, error)
	UpdateRolePermission(ctx context.Context, req *v1.UpdateRolePermissionRequest) error
}

func NewRoleService(
	service *Service,
	roleRepository repository.RoleRepository,
) RoleService {
	return &roleService{
		Service:        service,
		roleRepository: roleRepository,
	}
}

type roleService struct {
	*Service
	roleRepository repository.RoleRepository
}

func (s *roleService) ListRoles(ctx context.Context, req *v1.RoleSearchRequest) (*v1.RoleSearchResponseData, error) {
	list, total, err := s.roleRepository.ListRoles(ctx, req)
	if err != nil {
		return nil, err
	}
	data := &v1.RoleSearchResponseData{
		List:  make([]v1.RoleDataItem, 0),
		Total: total,
	}
	for _, role := range list {
		data.List = append(data.List, v1.RoleDataItem{
			ID:        role.ID,
			CreatedAt: role.CreatedAt.Format(constant.DateTimeLayout),
			UpdatedAt: role.UpdatedAt.Format(constant.DateTimeLayout),
			Name:      role.Name,
			Role:      role.Role,
		})

	}
	return data, nil
}

func (s *roleService) RoleCreate(ctx context.Context, req *v1.RoleRequest) error {
	_, err := s.roleRepository.GetRoleByCasbinRole(ctx, req.Role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.roleRepository.RoleCreate(ctx, &model.Role{
				Name: req.Name,
				Role: req.Role,
			})
		} else {
			return err
		}
	}
	return nil
}

func (s *roleService) RoleUpdate(ctx context.Context, id uint, req *v1.RoleRequest) error {
	return s.roleRepository.RoleUpdate(ctx, &model.Role{
		Model: gorm.Model{
			ID: id,
		},
		Name: req.Name,
	})
}

func (s *roleService) RoleDelete(ctx context.Context, id uint) error {
	old, err := s.roleRepository.GetRole(ctx, id)
	if err != nil {
		return err
	}
	if _, err := s.roleRepository.CasbinRoleDelete(ctx, old.Role); err != nil {
		return err
	}
	return s.roleRepository.RoleDelete(ctx, id)
}

func (s *roleService) GetRolePermission(ctx context.Context, role string) (*v1.GetRolePermissionResponseData, error) {
	// 获取角色对应的权限列表
	list, err := s.roleRepository.GetRolePermission(ctx, role)
	if err != nil {
		return nil, err
	}

	data := &v1.GetRolePermissionResponseData{
		List:  []string{},
		Total: 0,
	}
	for _, v := range list {
		if len(v) == 3 {
			data.List = append(data.List, strings.Join([]string{v[1], v[2]}, constant.PermSep))
		}
	}
	data.Total = int64(len(data.List))
	return data, nil
}

func (s *roleService) UpdateRolePermission(ctx context.Context, req *v1.UpdateRolePermissionRequest) error {
	permissions := map[string]struct{}{}
	for _, v := range req.List {
		perm := strings.Split(v, constant.PermSep)
		if len(perm) == 2 {
			permissions[v] = struct{}{}
		}

	}
	return s.roleRepository.UpdateRolePermission(ctx, req.Role, permissions)
}
