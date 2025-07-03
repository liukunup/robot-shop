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
	m.db.Migrator().DropTable(
		&model.User{},
		&model.Menu{},
		&model.Role{},
		&model.Api{},
		&model.Robot{},
	)
	if err := m.db.AutoMigrate(
		&model.User{},
		&model.Menu{},
		&model.Role{},
		&model.Api{},
		&model.Robot{},
	); err != nil {
		m.log.Error("user migrate error", zap.Error(err))
		return err
	}
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = m.db.Create(&model.User{
		Model:    gorm.Model{ID: 1},
		Username: "admin",
		Password: string(hashedPassword),
		Nickname: "超级管理员",
		Email:    "admin@example.com",
		Phone:    "12345678901",
	}).Error
	return m.db.Create(&model.User{
		Model:    gorm.Model{ID: 2},
		Username: "operator",
		Password: string(hashedPassword),
		Nickname: "运营人员",
		Email:    "operator@example.com",
		Phone:    "12345678901",
	}).Error
}

func (m *MigrateServer) initialRBAC(ctx context.Context) error {
	roles := []model.Role{
		{CasbinRole: constant.AdminRole, Name: "超级管理员"},
		{CasbinRole: "operator", Name: "运营人员"},
		{CasbinRole: "guest", Name: "访客"},
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
	_, err = m.e.AddRoleForUser(constant.AdminUserID, constant.AdminRole)
	if err != nil {
		m.log.Error("m.e.AddRoleForUser error", zap.Error(err))
		return err
	}
	menuList := make([]v1.MenuDataItem, 0)
	err = json.Unmarshal([]byte(menuData), &menuList)
	if err != nil {
		m.log.Error("json.Unmarshal error", zap.Error(err))
		return err
	}
	for _, item := range menuList {
		m.addPermissionForRole(constant.AdminRole, constant.MenuResourcePrefix+item.Path, "read")
	}
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
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/profile/basic", "read")
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/profile/advanced", "read")
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/profile", "read")
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/dashboard", "read")
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/dashboard/workplace", "read")
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/dashboard/analysis", "read")
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/account/settings", "read")
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/account/center", "read")
	m.addPermissionForRole(constant.OperatorRole, constant.MenuResourcePrefix+"/account", "read")
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

		// 基础API
		{Group: "基础API", Name: "登录", Path: "/v1/login", Method: http.MethodPost},
		{Group: "基础API", Name: "注册", Path: "/v1/register", Method: http.MethodPost},
		{Group: "基础API", Name: "获取当前用户信息", Path: "/v1/users/me", Method: http.MethodGet},
		{Group: "基础API", Name: "获取当前用户菜单", Path: "/v1/users/me/menus", Method: http.MethodGet},
		{Group: "基础API", Name: "获取当前用户权限", Path: "/v1/users/me/permissions", Method: http.MethodGet},

		// 用户管理
		{Group: "用户管理", Name: "获取用户列表", Path: "/v1/admin/users", Method: http.MethodGet},
		{Group: "用户管理", Name: "创建用户", Path: "/v1/admin/users", Method: http.MethodPost},
		{Group: "用户管理", Name: "更新用户", Path: "/v1/admin/users/:id", Method: http.MethodPut},
		{Group: "用户管理", Name: "删除用户", Path: "/v1/admin/users/:id", Method: http.MethodDelete},

		// 角色管理
		{Group: "角色管理", Name: "获取角色列表", Path: "/v1/admin/roles", Method: http.MethodGet},
		{Group: "角色管理", Name: "创建角色", Path: "/v1/admin/roles", Method: http.MethodPost},
		{Group: "角色管理", Name: "更新角色", Path: "/v1/admin/roles/:id", Method: http.MethodPut},
		{Group: "角色管理", Name: "删除角色", Path: "/v1/admin/roles/:id", Method: http.MethodDelete},
		{Group: "角色管理", Name: "获取角色权限", Path: "/v1/admin/roles/permissions", Method: http.MethodGet},
		{Group: "角色管理", Name: "更新角色权限", Path: "/v1/admin/roles/permissions", Method: http.MethodPut},

		// 菜单管理
		{Group: "菜单管理", Name: "获取菜单列表", Path: "/v1/admin/menus", Method: http.MethodGet},
		{Group: "菜单管理", Name: "创建菜单", Path: "/v1/admin/menus", Method: http.MethodPost},
		{Group: "菜单管理", Name: "更新菜单", Path: "/v1/admin/menus/:id", Method: http.MethodPut},
		{Group: "菜单管理", Name: "删除菜单", Path: "/v1/admin/menus/:id", Method: http.MethodDelete},

		// 接口管理
		{Group: "接口管理", Name: "获取接口列表", Path: "/v1/admin/apis", Method: http.MethodGet},
		{Group: "接口管理", Name: "创建接口", Path: "/v1/admin/apis", Method: http.MethodPost},
		{Group: "接口管理", Name: "更新接口", Path: "/v1/admin/apis/:id", Method: http.MethodPut},
		{Group: "接口管理", Name: "删除接口", Path: "/v1/admin/apis/:id", Method: http.MethodDelete},

		// 机器人管理
		{Group: "机器人管理", Name: "获取机器人列表", Path: "/v1/robots", Method: http.MethodGet},
		{Group: "机器人管理", Name: "获取机器人详情", Path: "/v1/robots/:id", Method: http.MethodGet},
		{Group: "机器人管理", Name: "创建机器人", Path: "/v1/robots", Method: http.MethodPost},
		{Group: "机器人管理", Name: "更新机器人", Path: "/v1/robots/:id", Method: http.MethodPut},
		{Group: "机器人管理", Name: "删除机器人", Path: "/v1/robots/:id", Method: http.MethodDelete},
	}

	return m.db.Create(&initialApis).Error
}

func (m *MigrateServer) initialMenuData(ctx context.Context) error {
	menuList := make([]v1.MenuDataItem, 0)
	err := json.Unmarshal([]byte(menuData), &menuList)
	if err != nil {
		m.log.Error("json.Unmarshal error", zap.Error(err))
		return err
	}
	menuListDb := make([]model.Menu, 0)
	for _, item := range menuList {
		menuListDb = append(menuListDb, model.Menu{
			Model: gorm.Model{
				ID: item.ID,
			},
			ParentID:  item.ParentID,
			Path:      item.Path,
			Redirect:  item.Redirect,
			Component: item.Component,
			Name:      item.Name,
			Icon:      item.Icon,
			Access:    item.Access,
			Weight:    item.Weight,
		})
	}
	return m.db.Create(&menuListDb).Error
}

var menuData = `[
  {
    "id": 1,
    "parentId": 0,
    "path": "/robot",
    "component": "./Robot",
    "name": "robot",
    "icon": "robot"
  },
  {
    "id": 1000,
    "parentId": 0,
    "path": "/admin",
    "name": "admin",
    "icon": "crown",
    "access": "admin"
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
    "component": "./Admin/User",
    "name": "user",
    "weight": 1
  },
  {
    "id": 1003,
    "parentId": 1000,
    "path": "/admin/role",
    "component": "./Admin/Role",
    "name": "role",
    "weight": 2
  },
  {
    "id": 1004,
    "parentId": 1000,
    "path": "/admin/menu",
    "component": "./Admin/Menu",
    "name": "menu",
    "weight": 3
  },
  {
    "id": 1005,
    "parentId": 1000,
    "path": "/admin/api",
    "component": "./Admin/Api",
    "name": "api",
    "weight": 4
  }
]`
