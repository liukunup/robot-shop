package repository

import (
	v1 "backend/api/v1"
	"backend/internal/constant"
	"backend/internal/model"
	"context"
	"fmt"
	"strings"
)

type RoleRepository interface {
	ListRoles(ctx context.Context, req *v1.GetRoleListRequest) ([]model.Role, int64, error)
	RoleUpdate(ctx context.Context, m *model.Role) error
	RoleCreate(ctx context.Context, m *model.Role) error
	RoleDelete(ctx context.Context, id uint) error

	GetRole(ctx context.Context, id uint) (model.Role, error)
	GetRoleBySid(ctx context.Context, sid string) (model.Role, error)

	CasbinRoleDelete(ctx context.Context, role string) error

	GetRolePermissions(ctx context.Context, role string) ([][]string, error)
	UpdateRolePermission(ctx context.Context, role string, permissions map[string]struct{}) error
}

func NewRoleRepository(
	repository *Repository,
) RoleRepository {
	return &roleRepository{
		Repository: repository,
	}
}

type roleRepository struct {
	*Repository
}

func (r *roleRepository) ListRoles(ctx context.Context, req *v1.GetRoleListRequest) ([]model.Role, int64, error) {
	var list []model.Role
	var total int64
	scope := r.DB(ctx).Model(&model.Role{})
	if req.Name != "" {
		scope = scope.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Sid != "" {
		scope = scope.Where("sid = ?", req.Sid)
	}
	if err := scope.Count(&total).Error; err != nil {
		return nil, total, err
	}
	if err := scope.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error; err != nil {
		return nil, total, err
	}
	return list, total, nil
}

func (r *roleRepository) RoleUpdate(ctx context.Context, m *model.Role) error {
	return r.DB(ctx).Where("id = ?", m.ID).UpdateColumn("name", m.Name).Error
}

func (r *roleRepository) RoleCreate(ctx context.Context, m *model.Role) error {
	return r.DB(ctx).Create(m).Error
}

func (r *roleRepository) RoleDelete(ctx context.Context, id uint) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&model.Role{}).Error
}

func (r *roleRepository) GetRole(ctx context.Context, id uint) (model.Role, error) {
	m := model.Role{}
	return m, r.DB(ctx).Where("id = ?", id).First(&m).Error
}

func (r *roleRepository) GetRoleBySid(ctx context.Context, sid string) (model.Role, error) {
	m := model.Role{}
	return m, r.DB(ctx).Where("sid = ?", sid).First(&m).Error
}

func (r *roleRepository) CasbinRoleDelete(ctx context.Context, role string) error {
	_, err := r.e.DeleteRole(role)
	return err
}

func (r *roleRepository) GetRolePermissions(ctx context.Context, role string) ([][]string, error) {
	return r.e.GetPermissionsForUser(role)
}

func (r *roleRepository) UpdateRolePermission(ctx context.Context, role string, newPermSet map[string]struct{}) error {
	if len(newPermSet) == 0 {
		return nil
	}
	// 获取当前角色的所有权限
	oldPermissions, err := r.e.GetPermissionsForUser(role)
	if err != nil {
		return err
	}

	// 将旧权限转换为 map 方便查找
	oldPermSet := make(map[string]struct{})
	for _, perm := range oldPermissions {
		if len(perm) == 3 {
			oldPermSet[strings.Join([]string{perm[1], perm[2]}, constant.PermSep)] = struct{}{}
		}
	}

	// 找出需要删除的权限
	var removePermissions [][]string
	for key, _ := range oldPermSet {
		if _, exists := newPermSet[key]; !exists {
			removePermissions = append(removePermissions, strings.Split(key, constant.PermSep))
		}
	}

	// 找出需要添加的权限
	var addPermissions [][]string
	for key, _ := range newPermSet {
		if _, exists := oldPermSet[key]; !exists {
			addPermissions = append(addPermissions, strings.Split(key, constant.PermSep))
		}

	}

	// 先移除多余的权限（使用 DeletePermissionForUser 逐条删除）
	for _, perm := range removePermissions {
		_, err := r.e.DeletePermissionForUser(role, perm...)
		if err != nil {
			return fmt.Errorf("移除权限失败: %v", err)
		}
	}

	// 再添加新的权限
	if len(addPermissions) > 0 {
		_, err = r.e.AddPermissionsForUser(role, addPermissions...)
		if err != nil {
			return fmt.Errorf("添加新权限失败: %v", err)
		}
	}

	return nil
}
