package service

import (
	v1 "backend/api/v1"
	"backend/internal/constant"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/email"
	"context"
	cryptoRand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	mathRand "math/rand"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	List(ctx context.Context, req *v1.UserSearchRequest) (*v1.UserSearchResponseData, error)
	Create(ctx context.Context, req *v1.UserRequest) error
	Update(ctx context.Context, uid uint, req *v1.UserRequest) error
	Delete(ctx context.Context, uid uint) error
	Get(ctx context.Context, uid uint) (*v1.UserDataItem, error)

	GetMenus(ctx context.Context, uid uint) (*v1.DynamicMenuResponseData, error)

	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, error)
	UpdatePassword(ctx context.Context, uid uint, req *v1.UpdatePasswordRequest) error
	ResetPassword(ctx context.Context, req *v1.ResetPasswordRequest) error
}

func NewUserService(
	service *Service,
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
	menuRepository repository.MenuRepository,
) UserService {
	return &userService{
		Service:        service,
		userRepository: userRepository,
		roleRepository: roleRepository,
		menuRepository: menuRepository,
	}
}

type userService struct {
	*Service
	userRepository repository.UserRepository
	roleRepository repository.RoleRepository
	menuRepository repository.MenuRepository
}

func (s *userService) List(ctx context.Context, req *v1.UserSearchRequest) (*v1.UserSearchResponseData, error) {
	// 获取用户列表
	list, total, err := s.userRepository.List(ctx, req)
	if err != nil {
		return nil, err
	}

	data := &v1.UserSearchResponseData{
		List:  make([]v1.UserDataItem, 0),
		Total: total,
	}
	for _, user := range list {
		// 获取用户角色
		roles, err := s.userRepository.GetRoles(ctx, user.ID)
		if err != nil {
			s.logger.Error("userRepository.GetUserRoles error", zap.Error(err))
			continue
		}
		// 转成角色对象
		roleList := make([]v1.RoleDataItem, 0)
		if len(roles) > 0 {
			for _, role := range roles {
				m, err2 := s.roleRepository.GetByCasbinRole(ctx, role)
				if err2 != nil {
					s.logger.Error("roleRepository.GetRoleByCasbinRole error", zap.Error(err2))
					continue
				}
				roleList = append(roleList, v1.RoleDataItem{
					ID:         m.ID,
					Name:       m.Name,
					CasbinRole: m.CasbinRole,
				})
			}
		}
		data.List = append(data.List, v1.UserDataItem{
			ID:        user.ID,
			CreatedAt: user.CreatedAt.Format(constant.DateTimeLayout),
			UpdatedAt: user.UpdatedAt.Format(constant.DateTimeLayout),
			Email:     user.Email,
			Username:  user.Username,
			Avatar:    user.Avatar,
			Nickname:  user.Nickname,
			Bio:       user.Bio,
			Language:  user.Language,
			Timezone:  user.Timezone,
			Theme:     user.Theme,
			Status:    user.Status,
			Roles:     roleList,
		})
	}

	return data, nil
}

func (s *userService) Create(ctx context.Context, req *v1.UserRequest) error {
	var err error

	// 检查邮箱是否已存在
	_, err = s.userRepository.GetByEmail(ctx, req.Email)
	if err == nil {
		return v1.ErrEmailAlreadyUse
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return v1.ErrInternalServerError
	}

	// 检查用户名是否已存在
	_, err = s.userRepository.GetByUsername(ctx, req.Username)
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
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword), // 用户通过邮箱激活账户并设置新密码
		Nickname: req.Nickname,
		Bio:      req.Bio,
		Language: req.Language,
		Timezone: req.Timezone,
		Theme:    req.Theme,
		Status:   req.Status,
	}
	// 创建用户
	if err = s.userRepository.Create(ctx, newUser); err != nil {
		return err
	}
	// 设置角色
	if err = s.userRepository.UpdateRoles(ctx, newUser.ID, req.Roles); err != nil {
		return err
	}
	return err
}

func (s *userService) Update(ctx context.Context, uid uint, req *v1.UserRequest) error {
	// 更新用户角色
	err := s.userRepository.UpdateRoles(ctx, uid, req.Roles)
	if err != nil {
		return err
	}
	// 更新用户
	data := map[string]interface{}{
		"email":    req.Email,
		"username": req.Username,
		"nickname": req.Nickname,
		"bio":      req.Bio,
		"language": req.Language,
		"timezone": req.Timezone,
		"theme":    req.Theme,
		"status":   req.Status,
	}
	return s.userRepository.Update(ctx, uid, data)
}

func (s *userService) Delete(ctx context.Context, uid uint) error {
	// 删除用户角色
	err := s.userRepository.DeleteRoles(ctx, uid)
	if err != nil {
		return err
	}
	// 删除用户
	return s.userRepository.Delete(ctx, uid)
}

func (s *userService) Get(ctx context.Context, uid uint) (*v1.UserDataItem, error) {
	// 获取用户
	user, err := s.userRepository.Get(ctx, uid)
	if err != nil {
		s.logger.WithContext(ctx).Error("userRepository.Get error", zap.Error(err))
		return nil, err
	}
	// 获取用户角色
	roles, err := s.userRepository.GetRoles(ctx, uid)
	if err != nil {
		s.logger.WithContext(ctx).Error("userRepository.GetRoles error", zap.Error(err))
		return nil, err
	}
	// 转成角色对象
	roleList := make([]v1.RoleDataItem, 0)
	if len(roles) > 0 {
		for _, role := range roles {
			m, err2 := s.roleRepository.GetByCasbinRole(ctx, role)
			if err2 != nil {
				s.logger.Error("roleRepository.GetRoleByCasbinRole error", zap.Error(err2))
				continue
			}
			roleList = append(roleList, v1.RoleDataItem{
				ID:         m.ID,
				Name:       m.Name,
				CasbinRole: m.CasbinRole,
			})
		}
	}
	return &v1.UserDataItem{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Format(constant.DateTimeLayout),
		UpdatedAt: user.UpdatedAt.Format(constant.DateTimeLayout),
		Email:     user.Email,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Avatar:    "https://cravatar.cn/avatar/245467ef31b6f0addc72b039b94122a4?s=100&f=y&r=g",
		Bio:       user.Bio,
		Language:  user.Language,
		Timezone:  user.Timezone,
		Theme:     user.Theme,
		Status:    user.Status,
		Roles:     roleList,
	}, nil
}

func (s *userService) GetMenus(ctx context.Context, uid uint) (*v1.DynamicMenuResponseData, error) {
	// 获取菜单列表
	menuList, err := s.menuRepository.ListAll(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Error("menuRepository.ListAll error", zap.Error(err))
		return nil, err
	}

	// 获取权限列表
	permList, err := s.userRepository.GetPermissions(ctx, uid)
	if err != nil {
		s.logger.WithContext(ctx).Error("userRepository.GetPermissions error", zap.Error(err))
		return nil, err
	}

	// 构建菜单权限映射
	permMap := map[string]struct{}{}
	for _, permission := range permList {
		if len(permission) == 3 && strings.HasPrefix(permission[1], constant.MenuResourcePrefix) {
			permMap[strings.TrimPrefix(permission[1], constant.MenuResourcePrefix)] = struct{}{}
		}
	}

	// -------------------- 构建动态菜单 --------------------
	// 第一轮遍历 构建ID到节点的映射
	menuMap := make(map[uint]*v1.MenuNode)
	for _, menu := range menuList {
		// 权限过滤
		if _, ok := permMap[menu.Path]; !ok {
			continue
		}
		// 构建节点
		menuNode := &v1.MenuNode{
			MenuDataItem: v1.MenuDataItem{
				ID:                 menu.ID,
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
			},
			Children: make([]*v1.MenuNode, 0),
		}
		menuMap[menu.ID] = menuNode
	}

	// 第二轮遍历 构建树结构
	menuRoot := make([]*v1.MenuNode, 0)
	for _, menu := range menuList {
		// 权限过滤
		if _, ok := permMap[menu.Path]; !ok {
			continue
		}

		menuNode := menuMap[menu.ID]

		if menu.ParentID == 0 { // 顶级菜单
			menuRoot = append(menuRoot, menuNode)
		} else { // 子菜单
			if parent, ok := menuMap[menu.ParentID]; ok {
				parent.Children = append(parent.Children, menuNode)
			}
		}
	}

	data := &v1.DynamicMenuResponseData{
		List: menuRoot,
	}
	return data, nil
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	// 检查邮箱是否已注册
	_, err := s.userRepository.GetByEmail(ctx, req.Email)
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
		if err = s.userRepository.Create(ctx, user); err != nil {
			return err
		}
		// TODO: other repo
		return nil
	})
	return err
}

func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (string, error) {
	// 查找指定用户
	user, err := s.userRepository.GetByUsernameOrEmail(ctx, req.Username, req.Username)
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

func (s *userService) UpdatePassword(ctx context.Context, uid uint, req *v1.UpdatePasswordRequest) error {
	user, err := s.userRepository.Get(ctx, uid)
	if err != nil {
		return err
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return v1.ErrUnauthorized
	}

	// 创建新密码哈希值
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	user.Password = string(hashedPassword)
	if err = s.userRepository.Update(ctx, uid, map[string]interface{}{
		"password": user.Password,
	}); err != nil {
		return err
	}

	return nil
}

func (s *userService) ResetPassword(ctx context.Context, req *v1.ResetPasswordRequest) error {
	// 查找指定用户
	user, err := s.userRepository.GetByEmail(ctx, req.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return v1.ErrUnauthorized
	}
	if err != nil {
		return v1.ErrInternalServerError
	}

	// 生成重置密码链接
	resetLink := fmt.Sprintf("https://robot-shop.com/reset-password?token=%d", user.ID)

	// 发送重置密码邮件
	if err = s.email.Send(&email.Message{
		To:      []string{user.Email},
		Subject: constant.ResetPasswordSubject,
		Text:    fmt.Sprintf(constant.ResetPasswordTextTemplate, user.Nickname, resetLink),
	}); err != nil {
		return err
	}

	return nil
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
