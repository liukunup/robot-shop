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
	ListRoles(ctx context.Context, req *v1.RoleSearchRequest) ([]model.Role, int64, error)
	RoleCreate(ctx context.Context, m *model.Role) error
	RoleUpdate(ctx context.Context, m *model.Role) error
	RoleDelete(ctx context.Context, id uint) error
	GetRole(ctx context.Context, id uint) (model.Role, error)

	GetRoleByCasbinRole(ctx context.Context, role string) (model.Role, error)

	GetRolePermission(ctx context.Context, role string) ([][]string, error)
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

func (r *roleRepository) ListRoles(ctx context.Context, req *v1.RoleSearchRequest) ([]model.Role, int64, error) {
	var list []model.Role
	var total int64
	scope := r.DB(ctx).Model(&model.Role{})
	if req.Name != "" {
		scope = scope.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Role != "" {
		scope = scope.Where("role = ?", req.Role)
	}
	if err := scope.Count(&total).Error; err != nil {
		return nil, total, err
	}
	if err := scope.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error; err != nil {
		return nil, total, err
	}
	return list, total, nil
}

func (r *roleRepository) RoleCreate(ctx context.Context, m *model.Role) error {
	return r.DB(ctx).Create(m).Error
}

func (r *roleRepository) RoleUpdate(ctx context.Context, m *model.Role) error {
	return r.DB(ctx).Model(&model.Role{}).Where("id = ?", m.ID).UpdateColumn("name", m.Name).Error
}

func (r *roleRepository) RoleDelete(ctx context.Context, id uint) error {
	return r.Transaction(ctx, func(ctx context.Context) error {
		db := r.DB(ctx)
		// 获取角色
		var role model.Role
		if err := db.Where("id = ?", id).First(&role).Error; err != nil {
			return err
		}
		// 删除角色对应的权限
		if _, err := r.e.DeleteRole(role.Role); err != nil {
			return err
		}
		// 删除角色
		if err := db.Where("id = ?", id).Delete(&model.Role{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *roleRepository) GetRole(ctx context.Context, id uint) (model.Role, error) {
	m := model.Role{}
	return m, r.DB(ctx).Where("id = ?", id).First(&m).Error
}

func (r *roleRepository) GetRoleByCasbinRole(ctx context.Context, role string) (model.Role, error) {
	m := model.Role{}
	return m, r.DB(ctx).Where("role = ?", role).First(&m).Error
}

func (r *roleRepository) GetRolePermission(ctx context.Context, role string) ([][]string, error) {
	return r.e.GetPermissionsForUser(role)
}

func (r *roleRepository) UpdateRolePermission(ctx context.Context, role string, newPermSet map[string]struct{}) error {
	// 如果没有新的权限需要更新
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
	var shouldRemovePermList [][]string
	for key, _ := range oldPermSet {
		if _, exists := newPermSet[key]; !exists {
			shouldRemovePermList = append(shouldRemovePermList, strings.Split(key, constant.PermSep))
		}
	}

	// 找出需要添加的权限
	var shouldAddPermList [][]string
	for key, _ := range newPermSet {
		if _, exists := oldPermSet[key]; !exists {
			shouldAddPermList = append(shouldAddPermList, strings.Split(key, constant.PermSep))
		}

	}

	// 先移除多余的权限（使用 DeletePermissionForUser 逐条删除）
	for _, perm := range shouldRemovePermList {
		_, err := r.e.DeletePermissionForUser(role, perm...)
		if err != nil {
			return fmt.Errorf("移除旧权限失败: %v", err)
		}
	}

	// 再添加新的权限
	if len(shouldAddPermList) > 0 {
		_, err = r.e.AddPermissionsForUser(role, shouldAddPermList...)
		if err != nil {
			return fmt.Errorf("添加新权限失败: %v", err)
		}
	}

	return nil
}
