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
	Get(ctx context.Context, id uint) (model.Role, error)
	List(ctx context.Context, req *v1.RoleSearchRequest) ([]model.Role, int64, error)
	Create(ctx context.Context, m *model.Role) error
	Update(ctx context.Context, m *model.Role) error
	Delete(ctx context.Context, id uint) error

	ListAll(ctx context.Context) ([]model.Role, error)

	GetByCasbinRole(ctx context.Context, casbinRole string) (model.Role, error)
	DeleteCasbinRole(ctx context.Context, casbinRole string) (bool, error)

	GetPermissions(ctx context.Context, casbinRole string) ([][]string, error)
	UpdatePermissions(ctx context.Context, casbinRole string, permissions map[string]struct{}) error
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

func (r *roleRepository) Get(ctx context.Context, id uint) (model.Role, error) {
	m := model.Role{}
	return m, r.DB(ctx).Where("id = ?", id).First(&m).Error
}

func (r *roleRepository) List(ctx context.Context, req *v1.RoleSearchRequest) ([]model.Role, int64, error) {
	var list []model.Role
	var total int64
	scope := r.DB(ctx).Model(&model.Role{})
	if req.Name != "" {
		scope = scope.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.CasbinRole != "" {
		scope = scope.Where("casbin_role = ?", req.CasbinRole)
	}
	if err := scope.Count(&total).Error; err != nil {
		return nil, total, err
	}
	if req.Page > 0 && req.PageSize > 0 {
		if err := scope.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error; err != nil {
			return nil, total, err
		}
	} else {
		if err := scope.Find(&list).Error; err != nil {
			return nil, total, err
		}
	}

	return list, total, nil
}

func (r *roleRepository) Create(ctx context.Context, m *model.Role) error {
	return r.DB(ctx).Create(m).Error
}

func (r *roleRepository) Update(ctx context.Context, m *model.Role) error {
	return r.DB(ctx).Model(&model.Role{}).Where("id = ?", m.ID).UpdateColumn("name", m.Name).Error
}

func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&model.Role{}).Error
}

func (r *roleRepository) ListAll(ctx context.Context) ([]model.Role, error) {
	var list []model.Role
	if err := r.DB(ctx).Model(&model.Role{}).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *roleRepository) GetByCasbinRole(ctx context.Context, casbinRole string) (model.Role, error) {
	m := model.Role{}
	return m, r.DB(ctx).Where("casbin_role = ?", casbinRole).First(&m).Error
}

func (r *roleRepository) DeleteCasbinRole(ctx context.Context, casbinRole string) (bool, error) {
	return r.e.DeleteRole(casbinRole)
}

func (r *roleRepository) GetPermissions(ctx context.Context, casbinRole string) ([][]string, error) {
	return r.e.GetPermissionsForUser(casbinRole)
}

func (r *roleRepository) UpdatePermissions(ctx context.Context, casbinRole string, newPermSet map[string]struct{}) error {
	// 如果没有新的权限需要更新
	if len(newPermSet) == 0 {
		return nil
	}

	// 获取当前角色的所有权限
	oldPermissions, err := r.e.GetPermissionsForUser(casbinRole)
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
		_, err := r.e.DeletePermissionForUser(casbinRole, perm...)
		if err != nil {
			return fmt.Errorf("移除旧权限失败: %v", err)
		}
	}

	// 再添加新的权限
	if len(shouldAddPermList) > 0 {
		_, err = r.e.AddPermissionsForUser(casbinRole, shouldAddPermList...)
		if err != nil {
			return fmt.Errorf("添加新权限失败: %v", err)
		}
	}

	return nil
}
