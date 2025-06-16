package service

import (
	v1 "backend/api/v1"
	"backend/internal/constant"
	"backend/internal/model"
	"backend/internal/repository"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	// Login
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, error)

	// User
	ListUsers(ctx context.Context, uid uint) (*v1.GetMenuResponseData, error)
	UserUpdate(ctx context.Context, req *v1.AdminUserUpdateRequest) error
	UserCreate(ctx context.Context, req *v1.AdminUserCreateRequest) error
	UserDelete(ctx context.Context, id uint) error
	GetUser(ctx context.Context, uid uint) (*v1.GetUserResponseData, error)

	// Menu
	ListMenus(ctx context.Context, uid uint) (*v1.GetMenuResponseData, error)
	MenuUpdate(ctx context.Context, req *v1.MenuUpdateRequest) error
	MenuCreate(ctx context.Context, req *v1.MenuCreateRequest) error
	MenuDelete(ctx context.Context, id uint) error
	GetAdminMenus(ctx context.Context) (*v1.GetMenuResponseData, error)

	// Role
	ListRoles(ctx context.Context, req *v1.GetRoleListRequest) (*v1.GetRolesResponseData, error)
	RoleUpdate(ctx context.Context, req *v1.RoleUpdateRequest) error
	RoleCreate(ctx context.Context, req *v1.RoleCreateRequest) error
	RoleDelete(ctx context.Context, id uint) error

	// API
	ListApis(ctx context.Context, req *v1.GetApisRequest) (*v1.GetApisResponseData, error)
	ApiUpdate(ctx context.Context, req *v1.ApiUpdateRequest) error
	ApiCreate(ctx context.Context, req *v1.ApiCreateRequest) error
	ApiDelete(ctx context.Context, id uint) error

	// Permission
	GetUserPermissions(ctx context.Context, uid uint) (*v1.GetUserPermissionsData, error)
	GetRolePermissions(ctx context.Context, role string) (*v1.GetRolePermissionsData, error)
	UpdateRolePermission(ctx context.Context, req *v1.UpdateRolePermissionRequest) error
}

func NewUserService(
	service *Service,
	userRepository repository.UserRepository,
	menuRepository repository.MenuRepository,
	roleRepository repository.RoleRepository,
	apiRepository repository.ApiRepository,
) UserService {
	return &userService{
		Service:        service,
		userRepository: userRepository,
		menuRepository: menuRepository,
		roleRepository: roleRepository,
		apiRepository:  apiRepository,
	}
}

type userService struct {
	*Service
	userRepository repository.UserRepository
	menuRepository repository.MenuRepository
	roleRepository repository.RoleRepository
	apiRepository  repository.ApiRepository
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	// check username
	user, err := s.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if err == nil && user != nil {
		return v1.ErrEmailAlreadyUse
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Generate user ID
	userId, err := s.sid.GenString()
	if err != nil {
		return err
	}
	user = &model.User{
		Username: "",
		Password: string(hashedPassword),
		Email:    req.Email,
	}
	// Transaction demo
	err = s.tm.Transaction(ctx, func(ctx context.Context) error {
		// Create a user
		if err = s.userRepository.UserCreate(ctx, user); err != nil {
			return err
		}
		// TODO: other repo
		return nil
	})
	return err
}

func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (string, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return "", v1.ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", err
	}
	token, err := s.jwt.GenToken(user.Id, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) ListUsers(ctx context.Context, req *v1.GetUsersRequest) (*v1.GetAdminUsersResponseData, error) {
	list, total, err := s.userRepository.ListUsers(ctx, req)
	if err != nil {
		return nil, err
	}
	data := &v1.GetAdminUsersResponseData{
		List:  make([]v1.UserDataItem, 0),
		Total: total,
	}
	for _, user := range list {
		roles, err := s.userRepository.GetUserRoles(ctx, user.ID)
		if err != nil {
			s.logger.Error("GetUserRoles error", zap.Error(err))
			continue
		}
		data.List = append(data.List, v1.UserDataItem{
			Email:     user.Email,
			ID:        user.ID,
			Nickname:  user.Nickname,
			Username:  user.Username,
			Phone:     user.Phone,
			Roles:     roles,
			CreatedAt: user.CreatedAt.Format(constant.DateTimeLayout),
			UpdatedAt: user.UpdatedAt.Format(constant.DateTimeLayout),
		})
	}
	return data, nil
}

func (s *userService) UserUpdate(ctx context.Context, req *v1.AdminUserUpdateRequest) error {
	old, _ := s.userRepository.GetUser(ctx, req.ID)
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		req.Password = string(hash)
	} else {
		req.Password = old.Password
	}
	err := s.userRepository.UpdateUserRoles(ctx, req.ID, req.Roles)
	if err != nil {
		return err
	}
	return s.userRepository.UserUpdate(ctx, &model.User{
		Model: gorm.Model{
			ID: req.ID,
		},
		Email:    req.Email,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Username: req.Username,
	})

}

func (s *userService) UserCreate(ctx context.Context, req *v1.AdminUserCreateRequest) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hash)
	err = s.userRepository.UserCreate(ctx, &model.User{
		Email:    req.Email,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return err
	}
	user, err := s.userRepository.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	err = s.userRepository.UpdateUserRoles(ctx, user.ID, req.Roles)
	if err != nil {
		return err
	}
	return err

}

func (s *userService) UserDelete(ctx context.Context, id uint) error {
	err := s.userRepository.DeleteUserRoles(ctx, id)
	if err != nil {
		return err
	}
	return s.userRepository.UserDelete(ctx, id)
}

func (s *userService) GetUser(ctx context.Context, uid uint) (*v1.GetUserResponseData, error) {
	user, err := s.userRepository.GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	roles, _ := s.userRepository.GetUserRoles(ctx, uid)

	return &v1.GetUserResponseData{
		Email:     user.Email,
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Phone:     user.Phone,
		Roles:     roles,
		CreatedAt: user.CreatedAt.Format(constant.DateTimeLayout),
		UpdatedAt: user.UpdatedAt.Format(constant.DateTimeLayout),
	}, nil
}

func (s *userService) GetMenus(ctx context.Context, uid uint) (*v1.GetMenuResponseData, error) {
	menuList, err := s.menuRepository.ListMenus(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Error("GetMenuList error", zap.Error(err))
		return nil, err
	}
	data := &v1.GetMenuResponseData{
		List: make([]v1.MenuDataItem, 0),
	}
	// 获取权限的菜单
	permissions, err := s.userRepository.GetUserPermissions(ctx, uid)
	if err != nil {
		return nil, err
	}
	menuPermMap := map[string]struct{}{}
	for _, permission := range permissions {
		// 防呆设置，超管可以看到所有菜单
		if convertor.ToString(uid) == constant.SuperAdminUserID {
			menuPermMap[strings.TrimPrefix(permission[1], constant.MenuResourcePrefix)] = struct{}{}
		} else {
			if len(permission) == 3 && strings.HasPrefix(permission[1], constant.MenuResourcePrefix) {
				menuPermMap[strings.TrimPrefix(permission[1], constant.MenuResourcePrefix)] = struct{}{}
			}
		}
	}

	for _, menu := range menuList {
		if _, ok := menuPermMap[menu.Path]; ok {
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
	}
	return data, nil
}

func (s *userService) MenuUpdate(ctx context.Context, req *v1.MenuUpdateRequest) error {
	return s.menuRepository.MenuUpdate(ctx, &model.Menu{
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
		Model: gorm.Model{
			ID: req.ID,
		},
	})
}

func (s *userService) MenuCreate(ctx context.Context, req *v1.MenuCreateRequest) error {
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

func (s *userService) MenuDelete(ctx context.Context, id uint) error {
	return s.menuRepository.MenuDelete(ctx, id)
}

func (s *userService) GetAdminMenus(ctx context.Context) (*v1.GetMenuResponseData, error) {
	menuList, err := s.menuRepository.ListMenus(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Error("GetMenuList error", zap.Error(err))
		return nil, err
	}
	data := &v1.GetMenuResponseData{
		List: make([]v1.MenuDataItem, 0),
	}
	for _, menu := range menuList {
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

func (s *userService) ListRoles(ctx context.Context, req *v1.GetRoleListRequest) (*v1.GetRolesResponseData, error) {
	list, total, err := s.roleRepository.ListRoles(ctx, req)
	if err != nil {
		return nil, err
	}
	data := &v1.GetRolesResponseData{
		List:  make([]v1.RoleDataItem, 0),
		Total: total,
	}
	for _, role := range list {
		data.List = append(data.List, v1.RoleDataItem{
			ID:        role.ID,
			Name:      role.Name,
			Sid:       role.Sid,
			UpdatedAt: role.UpdatedAt.Format(constant.DateTimeLayout),
			CreatedAt: role.CreatedAt.Format(constant.DateTimeLayout),
		})

	}
	return data, nil
}

func (s *userService) RoleUpdate(ctx context.Context, req *v1.RoleUpdateRequest) error {
	return s.roleRepository.RoleUpdate(ctx, &model.Role{
		Name: req.Name,
		Sid:  req.Sid,
		Model: gorm.Model{
			ID: req.ID,
		},
	})
}

func (s *userService) RoleCreate(ctx context.Context, req *v1.RoleCreateRequest) error {
	_, err := s.roleRepository.GetRoleBySid(ctx, req.Sid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.roleRepository.RoleCreate(ctx, &model.Role{
				Name: req.Name,
				Sid:  req.Sid,
			})
		} else {
			return err
		}
	}
	return nil
}

func (s *userService) RoleDelete(ctx context.Context, id uint) error {
	old, err := s.roleRepository.GetRole(ctx, id)
	if err != nil {
		return err
	}
	if err := s.roleRepository.CasbinRoleDelete(ctx, old.Sid); err != nil {
		return err
	}
	return s.roleRepository.RoleDelete(ctx, id)
}

func (s *userService) ListApis(ctx context.Context, req *v1.GetApisRequest) (*v1.GetApisResponseData, error) {
	list, total, err := s.apiRepository.GetApis(ctx, req)
	if err != nil {
		return nil, err
	}
	groups, err := s.apiRepository.GetApiGroups(ctx)
	if err != nil {
		return nil, err
	}
	data := &v1.GetApisResponseData{
		List:   make([]v1.ApiDataItem, 0),
		Total:  total,
		Groups: groups,
	}
	for _, api := range list {
		data.List = append(data.List, v1.ApiDataItem{
			CreatedAt: api.CreatedAt.Format(constant.DateTimeLayout),
			Group:     api.Group,
			ID:        api.ID,
			Method:    api.Method,
			Name:      api.Name,
			Path:      api.Path,
			UpdatedAt: api.UpdatedAt.Format(constant.DateTimeLayout),
		})
	}
	return data, nil
}

func (s *userService) ApiUpdate(ctx context.Context, req *v1.ApiUpdateRequest) error {
	return s.apiRepository.ApiUpdate(ctx, &model.Api{
		Group:  req.Group,
		Method: req.Method,
		Name:   req.Name,
		Path:   req.Path,
		Model: gorm.Model{
			ID: req.ID,
		},
	})
}

func (s *userService) ApiCreate(ctx context.Context, req *v1.ApiCreateRequest) error {
	return s.apiRepository.ApiCreate(ctx, &model.Api{
		Group:  req.Group,
		Method: req.Method,
		Name:   req.Name,
		Path:   req.Path,
	})
}

func (s *userService) ApiDelete(ctx context.Context, id uint) error {
	return s.apiRepository.ApiDelete(ctx, id)
}

func (s *userService) GetUserPermissions(ctx context.Context, uid uint) (*v1.GetUserPermissionsData, error) {
	data := &v1.GetUserPermissionsData{
		List: []string{},
	}
	list, err := s.userRepository.GetUserPermissions(ctx, uid)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if len(v) == 3 {
			data.List = append(data.List, strings.Join([]string{v[1], v[2]}, constant.PermSep))
		}
	}
	return data, nil
}
func (s *userService) GetRolePermissions(ctx context.Context, role string) (*v1.GetRolePermissionsData, error) {
	data := &v1.GetRolePermissionsData{
		List: []string{},
	}
	list, err := s.roleRepository.GetRolePermissions(ctx, role)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if len(v) == 3 {
			data.List = append(data.List, strings.Join([]string{v[1], v[2]}, constant.PermSep))
		}
	}
	return data, nil
}

func (s *userService) UpdateRolePermission(ctx context.Context, req *v1.UpdateRolePermissionRequest) error {
	permissions := map[string]struct{}{}
	for _, v := range req.List {
		perm := strings.Split(v, constant.PermSep)
		if len(perm) == 2 {
			permissions[v] = struct{}{}
		}

	}
	return s.roleRepository.UpdateRolePermission(ctx, req.Role, permissions)
}
