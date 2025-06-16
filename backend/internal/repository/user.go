package repository

import (
	v1 "backend/api/v1"
	"backend/internal/model"
	"context"

	"github.com/duke-git/lancet/v2/convertor"
	"go.uber.org/zap"
)

type UserRepository interface {
	ListUsers(ctx context.Context, req *v1.GetUsersRequest) ([]model.User, int64, error)
	UserUpdate(ctx context.Context, m *model.User) error
	UserCreate(ctx context.Context, m *model.User) error
	UserDelete(ctx context.Context, id uint) error

	GetUser(ctx context.Context, uid uint) (model.User, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)

	GetUserPermissions(ctx context.Context, uid uint) ([][]string, error)
	GetUserRoles(ctx context.Context, uid uint) ([]string, error)
	UpdateUserRoles(ctx context.Context, uid uint, roles []string) error
	DeleteUserRoles(ctx context.Context, uid uint) error
}

func NewUserRepository(
	repository *Repository,
) UserRepository {
	return &userRepository{
		Repository: repository,
	}
}

type userRepository struct {
	*Repository
}

func (r *userRepository) ListUsers(ctx context.Context, req *v1.GetUsersRequest) ([]model.User, int64, error) {
	var list []model.User
	var total int64
	scope := r.DB(ctx).Model(&model.User{})
	if req.Username != "" {
		scope = scope.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Nickname != "" {
		scope = scope.Where("nickname LIKE ?", "%"+req.Nickname+"%")
	}
	if req.Email != "" {
		scope = scope.Where("email LIKE ?", "%"+req.Email+"%")
	}
	if req.Phone != "" {
		scope = scope.Where("phone LIKE ?", "%"+req.Phone+"%")
	}
	if err := scope.Count(&total).Error; err != nil {
		return nil, total, err
	}
	if err := scope.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("id DESC").Find(&list).Error; err != nil {
		return nil, total, err
	}
	return list, total, nil
}

func (r *userRepository) UserUpdate(ctx context.Context, m *model.User) error {
	return r.DB(ctx).Where("id = ?", m.ID).Save(m).Error
}

func (r *userRepository) UserCreate(ctx context.Context, m *model.User) error {
	return r.DB(ctx).Create(m).Error
}

func (r *userRepository) UserDelete(ctx context.Context, id uint) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&model.User{}).Error
}

func (r *userRepository) GetUser(ctx context.Context, uid uint) (model.User, error) {
	m := model.User{}
	return m, r.DB(ctx).Where("id = ?", uid).First(&m).Error
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	m := model.User{}
	return m, r.DB(ctx).Where("username = ?", username).First(&m).Error
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	m := model.User{}
	return m, r.DB(ctx).Where("email = ?", email).First(&m).Error
}

func (r *userRepository) GetUserPermissions(ctx context.Context, uid uint) ([][]string, error) {
	return r.e.GetImplicitPermissionsForUser(convertor.ToString(uid))

}

func (r *userRepository) GetUserRoles(ctx context.Context, uid uint) ([]string, error) {
	return r.e.GetRolesForUser(convertor.ToString(uid))
}

func (r *userRepository) UpdateUserRoles(ctx context.Context, uid uint, roles []string) error {
	if len(roles) == 0 {
		_, err := r.e.DeleteRolesForUser(convertor.ToString(uid))
		return err
	}
	old, err := r.e.GetRolesForUser(convertor.ToString(uid))
	if err != nil {
		return err
	}
	oldMap := make(map[string]struct{})
	newMap := make(map[string]struct{})
	for _, v := range old {
		oldMap[v] = struct{}{}
	}
	for _, v := range roles {
		newMap[v] = struct{}{}
	}
	addRoles := make([]string, 0)
	delRoles := make([]string, 0)

	for key, _ := range oldMap {
		if _, exists := newMap[key]; !exists {
			delRoles = append(delRoles, key)
		}
	}
	for key, _ := range newMap {
		if _, exists := oldMap[key]; !exists {
			addRoles = append(addRoles, key)
		}
	}
	if len(addRoles) == 0 && len(delRoles) == 0 {
		return nil
	}
	for _, role := range delRoles {
		if _, err := r.e.DeleteRoleForUser(convertor.ToString(uid), role); err != nil {
			r.logger.WithContext(ctx).Error("DeleteRoleForUser error", zap.Error(err))
			return err
		}
	}

	_, err = r.e.AddRolesForUser(convertor.ToString(uid), addRoles)
	return err
}

func (r *userRepository) DeleteUserRoles(ctx context.Context, uid uint) error {
	_, err := r.e.DeleteRolesForUser(convertor.ToString(uid))
	return err
}
