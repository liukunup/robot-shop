package server

import (
	v1 "backend/api/v1"
	"backend/internal/constant"
	"backend/internal/model"
	"backend/pkg/log"
	"backend/pkg/sid"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MigrateServer struct {
	db  *gorm.DB
	log *log.Logger
	sid *sid.Sid
	e   *casbin.SyncedEnforcer
}

func NewMigrateServer(
	db *gorm.DB,
	log *log.Logger,
	sid *sid.Sid,
	e *casbin.SyncedEnforcer,
) *MigrateServer {
	return &MigrateServer{
		e:   e,
		db:  db,
		log: log,
		sid: sid,
	}
}

func (m *MigrateServer) Start(ctx context.Context) error {
	// 模型列表
	models := []interface{}{
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.Api{},
		&model.Robot{},
	}
	// 使用事务确保迁移原子性
	m.db.Transaction(func(tx *gorm.DB) error {
		for _, model := range models {
			// 增量迁移
			if err := tx.AutoMigrate(model); err != nil {
				name := reflect.TypeOf(model).Elem().Name()
				m.log.Error("Migration error",
					zap.String("model", name),
					zap.Error(err))
				return fmt.Errorf("failed to migrate %s: %w", name, err)
			}
		}

		return nil
	})

	err := m.initialUser(ctx)
	if err != nil {
		m.log.Error("initialUser error", zap.Error(err))
	}

	err = m.initialMenuData(ctx)
	if err != nil {
		m.log.Error("initialMenuData error", zap.Error(err))
	}

	err = m.initialApisData(ctx)
	if err != nil {
		m.log.Error("initialApisData error", zap.Error(err))
	}

	err = m.initialRBAC(ctx)
	if err != nil {
		m.log.Error("initialRBAC error", zap.Error(err))
	}

	m.log.Info("AutoMigrate success")
	os.Exit(0)
	return nil
}
func (m *MigrateServer) Stop(ctx context.Context) error {
	m.log.Info("AutoMigrate stop")
	return nil
}

func (m *MigrateServer) initialUser(ctx context.Context) error {
	// 注意：在生产环境中请务必使用强密码！！！
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		m.log.Fatal("bcrypt.GenerateFromPassword error", zap.Error(err))
		return err
	}

	// 创建 超级管理员
	if err = m.db.Create(&model.User{
		Model:    gorm.Model{ID: 1},
		Username: "admin",
		Password: string(hashedPassword),
		Avatar:   "https://cravatar.cn/avatar/245467ef31b6f0addc72b039b94122a4?s=100&f=y&r=g",
		Nickname: "超级管理员",
		Email:    "admin@example.com",
		Status:   1,
	}).Error; err != nil {
		return err
	}

	// 创建 运营人员
	if err = m.db.Create(&model.User{
		Model:    gorm.Model{ID: 2},
		Username: "operator",
		Password: string(hashedPassword),
		Avatar:   "https://cravatar.cn/avatar/hash?s=100&d=robohash",
		Nickname: "运营人员",
		Email:    "operator@example.com",
		Status:   1,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (m *MigrateServer) initialRBAC(ctx context.Context) error {
	// 创建角色
	roles := []model.Role{
		{CasbinRole: constant.AdminRole, Name: "超级管理员"},
		{CasbinRole: constant.OperatorRole, Name: "运营人员"},
		{CasbinRole: constant.UserRole, Name: "普通用户"},
		{CasbinRole: constant.GuestRole, Name: "游客"},
	}
	if err := m.db.Create(&roles).Error; err != nil {
		return err
	}

	m.e.ClearPolicy()
	err := m.e.SavePolicy()
	if err != nil {
		m.log.Error("m.e.SavePolicy error", zap.Error(err))
		return err
	}

	// 超级管理员 + 用户
	_, err = m.e.AddRoleForUser(constant.AdminUserID, constant.AdminRole)
	if err != nil {
		m.log.Error("m.e.AddRoleForUser error", zap.Error(err))
		return err
	}
	// 超级管理员 + 菜单
	menuList := make([]model.Menu, 0)
	err = m.db.Find(&menuList).Error
	if err != nil {
		m.log.Error("m.db.Find(&menuList).Error error", zap.Error(err))
		return err
	}
	for _, menu := range menuList {
		m.addPermissionForRole(constant.AdminRole, constant.MenuResourcePrefix+menu.Path, constant.PermRead)
	}
	// 超级管理员 + 接口
	apiList := make([]model.Api, 0)
	err = m.db.Find(&apiList).Error
	if err != nil {
		m.log.Error("m.db.Find(&apiList).Error error", zap.Error(err))
		return err
	}
	for _, api := range apiList {
		m.addPermissionForRole(constant.AdminRole, constant.ApiResourcePrefix+api.Path, api.Method)
	}

	// 添加运营人员权限
	_, err = m.e.AddRoleForUser("2", constant.OperatorRole)
	if err != nil {
		m.log.Error("m.e.AddRoleForUser error", zap.Error(err))
		return err
	}
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/profile/basic", constant.PermRead)
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/profile/advanced", constant.PermRead)
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/profile", constant.PermRead)
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/dashboard", constant.PermRead)
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/dashboard/workplace", constant.PermRead)
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/dashboard/analysis", constant.PermRead)
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/account/settings", constant.PermRead)
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/account/center", constant.PermRead)
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/account", constant.PermRead)
	m.addPermissionForRole(constant.OperatorRole, constant.ApiResourcePrefix+"/v1/menus", http.MethodGet)
	m.addPermissionForRole(constant.OperatorRole, constant.ApiResourcePrefix+"/v1/admin/user", http.MethodGet)

	return nil
}

func (m *MigrateServer) addPermissionForRole(role, resource, action string) {
	_, err := m.e.AddPermissionForUser(role, resource, action)
	if err != nil {
		m.log.Sugar().Info("为角色 %s 添加权限 %s:%s 失败: %v", role, resource, action, err)
		return
	}
	fmt.Printf("为角色 %s 添加权限: %s %s\n", role, resource, action)
}

func (m *MigrateServer) initialApisData(ctx context.Context) error {
	initialApis := []model.Api{

		// 基础接口
		{Group: "Base", Name: "注册", Path: "/v1/register", Method: http.MethodPost},
		{Group: "Base", Name: "登录", Path: "/v1/login", Method: http.MethodPost},
		{Group: "Base", Name: "重置密码", Path: "/v1/reset-password", Method: http.MethodPost},
		{Group: "Base", Name: "刷新令牌", Path: "/v1/refresh-token", Method: http.MethodPost},

		// 用户
		{Group: "User", Name: "获取指定用户", Path: "/v1/users/:id", Method: http.MethodGet},

		// 用户
		{Group: "User", Name: "获取当前用户", Path: "/v1/users/profile", Method: http.MethodGet},
		{Group: "User", Name: "更新当前用户", Path: "/v1/users/profile", Method: http.MethodPut},
		{Group: "User", Name: "更新当前用户的头像", Path: "/v1/users/profile/avatar", Method: http.MethodPut},
		{Group: "User", Name: "获取当前用户的菜单", Path: "/v1/users/menus", Method: http.MethodGet},
		{Group: "User", Name: "更新当前用户的密码", Path: "/v1/users/password", Method: http.MethodPut},

		// 用户管理
		{Group: "User", Name: "获取用户列表", Path: "/v1/admin/users", Method: http.MethodGet},
		{Group: "User", Name: "创建用户", Path: "/v1/admin/users", Method: http.MethodPost},
		{Group: "User", Name: "更新用户", Path: "/v1/admin/users/:id", Method: http.MethodPut},
		{Group: "User", Name: "删除用户", Path: "/v1/admin/users/:id", Method: http.MethodDelete},

		// 角色管理
		{Group: "Role", Name: "获取角色列表", Path: "/v1/admin/roles", Method: http.MethodGet},
		{Group: "Role", Name: "创建角色", Path: "/v1/admin/roles", Method: http.MethodPost},
		{Group: "Role", Name: "更新角色", Path: "/v1/admin/roles/:id", Method: http.MethodPut},
		{Group: "Role", Name: "删除角色", Path: "/v1/admin/roles/:id", Method: http.MethodDelete},
		{Group: "Role", Name: "获取全部角色", Path: "/v1/admin/roles/all", Method: http.MethodGet},
		{Group: "Role", Name: "获取角色权限", Path: "/v1/admin/roles/permissions", Method: http.MethodGet},
		{Group: "Role", Name: "更新角色权限", Path: "/v1/admin/roles/permissions", Method: http.MethodPut},

		// 菜单管理
		{Group: "Menu", Name: "获取菜单列表", Path: "/v1/admin/menus", Method: http.MethodGet},
		{Group: "Menu", Name: "创建菜单", Path: "/v1/admin/menus", Method: http.MethodPost},
		{Group: "Menu", Name: "更新菜单", Path: "/v1/admin/menus/:id", Method: http.MethodPut},
		{Group: "Menu", Name: "删除菜单", Path: "/v1/admin/menus/:id", Method: http.MethodDelete},

		// 接口管理
		{Group: "API", Name: "获取接口列表", Path: "/v1/admin/apis", Method: http.MethodGet},
		{Group: "API", Name: "创建接口", Path: "/v1/admin/apis", Method: http.MethodPost},
		{Group: "API", Name: "更新接口", Path: "/v1/admin/apis/:id", Method: http.MethodPut},
		{Group: "API", Name: "删除接口", Path: "/v1/admin/apis/:id", Method: http.MethodDelete},

		// 机器人
		{Group: "Robot", Name: "获取机器人列表", Path: "/v1/robots", Method: http.MethodGet},
		{Group: "Robot", Name: "获取机器人详情", Path: "/v1/robots/:id", Method: http.MethodGet},
		{Group: "Robot", Name: "创建机器人", Path: "/v1/robots", Method: http.MethodPost},
		{Group: "Robot", Name: "更新机器人", Path: "/v1/robots/:id", Method: http.MethodPut},
		{Group: "Robot", Name: "删除机器人", Path: "/v1/robots/:id", Method: http.MethodDelete},
	}

	return m.db.Create(&initialApis).Error
}

func (m *MigrateServer) initialMenuData(ctx context.Context) error {
	menuList := make([]v1.MenuDataItem, 0)
	if err := json.Unmarshal([]byte(menuData), &menuList); err != nil {
		m.log.Error("json.Unmarshal error", zap.Error(err))
		return err
	}
	menuListDb := make([]model.Menu, 0)
	for _, item := range menuList {
		menuListDb = append(menuListDb, model.Menu{
			Model: gorm.Model{
				ID: item.ID,
			},
			ParentID:           item.ParentID,
			Icon:               item.Icon,
			Name:               item.Name,
			Path:               item.Path,
			Component:          item.Component,
			Access:             item.Access,
			Redirect:           item.Redirect,
			Target:             item.Target,
			HideChildrenInMenu: item.HideChildrenInMenu,
			HideInMenu:         item.HideInMenu,
			FlatMenu:           item.FlatMenu,
			Disabled:           item.Disabled,
			Tooltip:            item.Tooltip,
			DisabledTooltip:    item.DisabledTooltip,
			Key:                item.Key,
			ParentKeys:         item.ParentKeys,
		})
	}
	return m.db.Create(&menuListDb).Error
}

var menuData = `[
  {
    "id": 1,
    "path": "/",
	"redirect": "/welcome"
  },
  {
    "id": 2,
    "path": "/welcome",
    "name": "welcome",
    "icon": "smile",
	"component": "@/pages/Welcome"
  },
  {
    "id": 3,
    "path": "/robot",
    "name": "robot",
    "icon": "robot",
	"component": "@/pages/Robot",
	"access": "canUser"
  },
  {
    "id": 999,
    "path": "/profile",
    "name": "profile",
    "icon": "profile",
	"component": "@/pages/Profile",
	"access": "canUser"
  },
  {
    "id": 1000,
    "path": "/admin",
    "name": "admin",
    "icon": "crown",
    "access": "canAdmin"
  },
  {
    "id": 1001,
    "parentId": 1000,
    "path": "/admin",
	"redirect": "/admin/user"
  },
  {
    "id": 1002,
    "parentId": 1000,
    "path": "/admin/user",
    "name": "user",
	"component": "@/pages/Admin/User"
  },
  {
    "id": 1003,
    "parentId": 1000,
    "path": "/admin/role",
    "name": "role",
	"component": "@/pages/Admin/Role"
  },
  {
    "id": 1004,
    "parentId": 1000,
    "path": "/admin/menu",
    "name": "menu",
	"component": "@/pages/Admin/Menu"
  },
  {
    "id": 1005,
    "parentId": 1000,
    "path": "/admin/api",
    "name": "api",
	"component": "@/pages/Admin/Api"
  }
]`
