package service

import (
	v1 "backend/api/v1"
	"backend/internal/constant"
	"backend/internal/model"
	"backend/internal/repository"
	"context"
	cryptoRand "crypto/rand"
	"encoding/base64"
	"errors"
	mathRand "math/rand"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	// CRUD
	ListUsers(ctx context.Context, req *v1.UserSearchRequest) (*v1.UserSearchResponseData, error)
	UserCreate(ctx context.Context, req *v1.UserRequest) error
	UserUpdate(ctx context.Context, uid uint, req *v1.UserRequest) error
	UserDelete(ctx context.Context, uid uint) error
	GetUser(ctx context.Context, uid uint) (*v1.UserDataItem, error)
	// ATTR
	GetUserMenu(ctx context.Context, uid uint) (*v1.MenuSearchResponseData, error)
	GetUserPermission(ctx context.Context, uid uint) (*v1.UserPermissionResponseData, error)
	// Management
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, error)
}

func NewUserService(
	service *Service,
	userRepository repository.UserRepository,
	menuRepository repository.MenuRepository,
) UserService {
	return &userService{
		Service:        service,
		userRepository: userRepository,
		menuRepository: menuRepository,
	}
}

type userService struct {
	*Service
	userRepository repository.UserRepository
	menuRepository repository.MenuRepository
}

func (s *userService) ListUsers(ctx context.Context, req *v1.UserSearchRequest) (*v1.UserSearchResponseData, error) {
	// 获取用户列表
	list, total, err := s.userRepository.ListUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	data := &v1.UserSearchResponseData{
		List:  make([]v1.UserDataItem, 0),
		Total: total,
	}
	for _, user := range list {
		// 获取用户角色
		roles, err := s.userRepository.GetUserRoles(ctx, user.ID)
		if err != nil {
			s.logger.Error("userRepository.GetRoles error", zap.Error(err))
			continue
		}

		data.List = append(data.List, v1.UserDataItem{
			ID:        user.ID,
			CreatedAt: user.CreatedAt.Format(constant.DateTimeLayout),
			UpdatedAt: user.UpdatedAt.Format(constant.DateTimeLayout),
			Username:  user.Username,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Email:     user.Email,
			Phone:     user.Phone,
			Status:    user.Status,
			Roles:     roles,
		})
	}

	return data, nil
}

func (s *userService) UserCreate(ctx context.Context, req *v1.UserRequest) error {
	var err error

	// 检查邮箱是否已存在
	_, err = s.userRepository.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return v1.ErrEmailAlreadyUse
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return v1.ErrInternalServerError
	}

	// 检查用户名是否已存在
	_, err = s.userRepository.GetUserByUsername(ctx, req.Username)
	if err == nil {
		return v1.ErrUsernameAlreadyUse
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return v1.ErrInternalServerError
	}

	// 使用随机生成的密码
	randomPassword := generateRandomPassword(16)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 使用随机生成的昵称
	if req.Nickname == "" {
		req.Nickname = generateHumanNickname()
	}

	// 构造新用户对象
	newUser := &model.User{
		Username: req.Username,
		Password: string(hashedPassword), // 用户通过邮箱激活账户并设置新密码
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
	}
	err = s.tm.Transaction(ctx, func(ctx context.Context) error {
		// 创建用户
		if err = s.userRepository.UserCreate(ctx, newUser); err != nil {
			return err
		}
		// 设置角色
		if err = s.userRepository.UpdateUserRoles(ctx, newUser.ID, req.Roles); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *userService) UserUpdate(ctx context.Context, uid uint, req *v1.UserRequest) error {
	// 更新用户角色
	err := s.userRepository.UpdateUserRoles(ctx, uid, req.Roles)
	if err != nil {
		return err
	}
	// 更新用户
	data := map[string]interface{}{
		"username": req.Username,
		"nickname": req.Nickname,
		"email":    req.Email,
		"phone":    req.Phone,
		"status":   req.Status,
	}
	return s.userRepository.UserUpdate(ctx, uid, data)
}

func (s *userService) UserDelete(ctx context.Context, uid uint) error {
	// 删除用户角色
	err := s.userRepository.DeleteUserRoles(ctx, uid)
	if err != nil {
		return err
	}
	// 删除用户
	return s.userRepository.UserDelete(ctx, uid)
}

func (s *userService) GetUser(ctx context.Context, uid uint) (*v1.UserDataItem, error) {
	// 获取用户
	user, err := s.userRepository.GetUser(ctx, uid)
	if err != nil {
		s.logger.WithContext(ctx).Error("userRepository.Get error", zap.Error(err))
		return nil, err
	}
	// 获取用户角色
	roles, err := s.userRepository.GetUserRoles(ctx, uid)
	if err != nil {
		s.logger.WithContext(ctx).Error("userRepository.GetRoles error", zap.Error(err))
		return nil, err
	}

	return &v1.UserDataItem{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Format(constant.DateTimeLayout),
		UpdatedAt: user.UpdatedAt.Format(constant.DateTimeLayout),
		Username:  user.Username,
		Nickname:  user.Nickname,
		Avatar:    "https://cravatar.cn/avatar/245467ef31b6f0addc72b039b94122a4?s=100&f=y&r=g",
		Email:     user.Email,
		Phone:     user.Phone,
		Status:    user.Status,
		Roles:     roles,
	}, nil
}

func (s *userService) GetUserMenu(ctx context.Context, uid uint) (*v1.MenuSearchResponseData, error) {
	// 获取菜单列表
	menuList, total, err := s.menuRepository.ListMenus(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Error("menuRepository.List error", zap.Error(err))
		return nil, err
	}
	data := &v1.MenuSearchResponseData{
		List:  make([]v1.MenuDataItem, 0),
		Total: total,
	}

	// 获取权限列表
	permList, err := s.userRepository.GetUserPermission(ctx, uid)
	if err != nil {
		s.logger.WithContext(ctx).Error("userRepository.GetPermissions error", zap.Error(err))
		return nil, err
	}
	// 构建菜单权限映射
	menuPermMap := map[string]struct{}{}
	// 超管可以看到所有菜单
	if convertor.ToString(uid) == constant.AdminUserID {
		for _, permission := range permList {
			menuPermMap[strings.TrimPrefix(permission[1], constant.MenuResourcePrefix)] = struct{}{}
		}
	} else {
		for _, permission := range permList {
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

func (s *userService) GetUserPermission(ctx context.Context, uid uint) (*v1.UserPermissionResponseData, error) {
	// 获取权限列表
	permList, err := s.userRepository.GetUserPermission(ctx, uid)
	if err != nil {
		s.logger.WithContext(ctx).Error("userRepository.GetPermissions error", zap.Error(err))
		return nil, err
	}

	data := &v1.UserPermissionResponseData{
		List: []string{},
	}
	for _, v := range permList {
		if len(v) == 3 {
			data.List = append(data.List, strings.Join([]string{v[1], v[2]}, constant.PermSep))
		}
	}
	return data, nil
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	// 检查邮箱是否已注册
	_, err := s.userRepository.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return v1.ErrEmailAlreadyUse
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return v1.ErrInternalServerError
	}

	// 创建密码哈希值
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 从邮箱中提取默认用户名
	parts := strings.Split(req.Email, "@")
	if len(parts) != 2 {
		return v1.ErrInternalServerError
	}
	defaultUsername := parts[0]

	// 生成默认昵称
	nickname := generateHumanNickname()

	// 构造新用户对象
	user := &model.User{
		Username: defaultUsername,
		Password: string(hashedPassword),
		Nickname: nickname,
		Email:    req.Email,
		Status:   0,
	}
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
	// 查找指定用户
	user, err := s.userRepository.GetUserByUsernameOrEmail(ctx, req.Username, req.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", v1.ErrUnauthorized
	}
	if err != nil {
		return "", v1.ErrInternalServerError
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", err
	}

	// 创建AccessToken
	token, err := s.jwt.GenToken(user.ID, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return "", err
	}

	return token, nil
}

// 生成指定长度的随机密码
func generateRandomPassword(length int) string {
	b := make([]byte, length)
	_, err := cryptoRand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

// 生成一个人性化的昵称
func generateHumanNickname() string {

	adjectives := []string{
		"阳光的", "温柔的", "睿智的", "活泼的", "优雅的",
		"勇敢的", "幽默的", "神秘的", "开朗的", "沉稳的",
		"可爱的", "聪明的", "热情的", "冷静的", "浪漫的",
		"乐观的", "坚强的", "细心的", "真诚的", "大方的",
		"自由的", "独特的", "时尚的", "古典的", "现代的",
		"快乐的", "宁静的", "梦幻的", "激情的", "稳重的",
	}

	nouns := []string{
		"小明", "小华", "子轩", "雨桐", "浩然",
		"诗涵", "宇航", "欣怡", "俊杰", "雅婷",
		"志强", "美玲", "文博", "雪梅", "家豪",
		"丽娜", "建国", "婷婷", "海涛", "静怡",
		"小龙", "佳琪", "宏伟", "芳芳", "志明",
		"小雨", "天宇", "思琪", "大伟", "梦瑶",
	}

	seededRand := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	return adjectives[seededRand.Intn(len(adjectives))] + nouns[seededRand.Intn(len(nouns))]
}
