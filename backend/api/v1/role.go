package v1

// for Search
type RoleSearchRequest struct {
	Page     int    `form:"page" binding:"required" example:"1"`      // 页码
	PageSize int    `form:"pageSize" binding:"required" example:"10"` // 分页大小
	Name     string `form:"name" example:"admin"`                     // 筛选项: 角色名 模糊匹配
	Role     string `form:"role" example:"1"`                         // 筛选项: Role 精确匹配
}
type RoleDataItem struct {
	ID        uint   `json:"id"`                  // ID
	CreatedAt string `json:"createdAt,omitempty"` // 创建时间
	UpdatedAt string `json:"updatedAt,omitempty"` // 更新时间
	Name      string `json:"name"`                // 角色名
	Role      string `json:"role"`                // Casbin Role
} // @name Role
type RoleSearchResponseData struct {
	List  []RoleDataItem `json:"list"`
	Total int64          `json:"total"`
} // @name RoleList
type RoleSearchResponse struct {
	Response
	Data RoleSearchResponseData
}

// for Get
type RoleResponse struct {
	Response
	Data RoleDataItem
}

// for Create | Update
type RoleRequest struct {
	Name string `json:"name" binding:"required" example:"admin"` // 角色名
	Role string `json:"role" binding:"required" example:"1"`     // Casbin Role
}

// Role Permission
type GetRolePermissionRequest struct {
	Role string `json:"role" binding:"required" example:"admin"` // 角色名
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
	Role string   `form:"role" binding:"required" example:"admin"` // 角色名
	List []string `form:"list" binding:"required" example:""`      // 权限列表
}
