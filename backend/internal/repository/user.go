package repository

import (
	v1 "backend/api/v1"
	"backend/internal/model"
	"context"

	"github.com/duke-git/lancet/v2/convertor"
	"go.uber.org/zap"
)

type UserRepository interface {
	Get(ctx context.Context, uid uint) (model.User, error)
	List(ctx context.Context, req *v1.UserSearchRequest) ([]model.User, int64, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, uid uint, data map[string]interface{}) error
	Delete(ctx context.Context, uid uint) error

	GetByUsername(ctx context.Context, username string) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	GetByUsernameOrEmail(ctx context.Context, username string, email string) (model.User, error)

	GetPermissions(ctx context.Context, uid uint) ([][]string, error)
	GetRoles(ctx context.Context, uid uint) ([]string, error)
	UpdateRoles(ctx context.Context, uid uint, roles []string) error
	DeleteRoles(ctx context.Context, uid uint) error
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

func (r *userRepository) Get(ctx context.Context, uid uint) (model.User, error) {
	m := model.User{}
	return m, r.DB(ctx).Where("id = ?", uid).First(&m).Error
}

func (r *userRepository) List(ctx context.Context, req *v1.UserSearchRequest) ([]model.User, int64, error) {
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
	if err := scope.Count(&total).Error; err != nil {
		return nil, total, err
	}
	if err := scope.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("id ASC").Find(&list).Error; err != nil {
		return nil, total, err
	}
	return list, total, nil
}

func (r *userRepository) Create(ctx context.Context, m *model.User) error {
	return r.DB(ctx).Create(m).Error
}

func (r *userRepository) Update(ctx context.Context, uid uint, data map[string]interface{}) error {
	return r.DB(ctx).Model(&model.User{}).Where("id = ?", uid).Updates(data).Error
}

func (r *userRepository) Delete(ctx context.Context, uid uint) error {
	return r.DB(ctx).Where("id = ?", uid).Delete(&model.User{}).Error
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	m := model.User{}
	return m, r.DB(ctx).Where("username = ?", username).First(&m).Error
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	m := model.User{}
	return m, r.DB(ctx).Where("email = ?", email).First(&m).Error
}

func (r *userRepository) GetByUsernameOrEmail(ctx context.Context, username string, email string) (model.User, error) {
	m := model.User{}
	return m, r.DB(ctx).Where("username = ? OR email = ?", username, email).First(&m).Error
}

func (r *userRepository) GetPermissions(ctx context.Context, uid uint) ([][]string, error) {
	return r.e.GetImplicitPermissionsForUser(convertor.ToString(uid))
}

func (r *userRepository) GetRoles(ctx context.Context, uid uint) ([]string, error) {
	return r.e.GetRolesForUser(convertor.ToString(uid))
}

func (r *userRepository) UpdateRoles(ctx context.Context, uid uint, roles []string) error {
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

func (r *userRepository) DeleteRoles(ctx context.Context, uid uint) error {
	_, err := r.e.DeleteRolesForUser(convertor.ToString(uid))
	return err
}
