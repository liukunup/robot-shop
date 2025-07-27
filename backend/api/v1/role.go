package v1

// CRUD
type RoleSearchRequest struct {
	Page       int    `form:"page" binding:"omitempty,min=1" example:"1"`              // 页码
	PageSize   int    `form:"pageSize" binding:"omitempty,min=1,max=100" example:"10"` // 分页大小
	Name       string `form:"name" example:"Admin"`                                    // 筛选项: 角色名 模糊匹配
	CasbinRole string `form:"casbinRole" example:"admin"`                              // 筛选项: Casbin-Role 精确匹配
}
type RoleDataItem struct {
	ID         uint   `json:"id,omitempty"`                                       // ID
	CreatedAt  string `json:"createdAt,omitempty"  example:"2006-01-02 15:04:05"` // 创建时间
	UpdatedAt  string `json:"updatedAt,omitempty"  example:"2006-01-02 15:04:05"` // 更新时间
	Name       string `json:"name" example:"Admin"`                               // 角色名
	CasbinRole string `json:"casbinRole" example:"admin"`                         // Casbin-Role
} // @name Role
type RoleSearchResponseData struct {
	List  []RoleDataItem `json:"list"`  // 列表
	Total int64          `json:"total"` // 总数
} // @name RoleList
type RoleSearchResponse struct {
	Response
	Data RoleSearchResponseData
}

type RoleResponse struct {
	Response
	Data RoleDataItem
}

type RoleRequest struct {
	Name       string `json:"name" binding:"required" example:"Admin"`       // 角色名
	CasbinRole string `json:"casbinRole" binding:"required" example:"admin"` // Casbin-Role
}

// Permission Management
type GetRolePermissionRequest struct {
	CasbinRole string `json:"casbinRole" binding:"required" example:"admin"` // Casbin-Role
}
type GetRolePermissionResponseData struct {
	List  []string `json:"list"`  // 列表
	Total int64    `json:"total"` // 总数
}
type GetRolePermissionResponse struct {
	Response
	Data GetRolePermissionResponseData
}
type UpdateRolePermissionRequest struct {
	CasbinRole string   `json:"casbinRole" binding:"required" example:"admin"` // Casbin-Role
	List       []string `form:"list" binding:"required"`                       // 权限列表
}
