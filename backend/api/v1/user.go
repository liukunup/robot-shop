package v1

// for Search
type UserSearchRequest struct {
	Page     int    `form:"page" binding:"required" example:"1"`             // 页码
	PageSize int    `form:"pageSize" binding:"required" example:"10"`        // 分页大小
	Username string `form:"username" binding:"" example:"zhangsan"`          // 用户名
	Nickname string `form:"nickname" binding:"" example:"Jackal"`            // 昵称
	Phone    string `form:"phone" binding:"" example:"+86-13966668888"`      // 手机
	Email    string `form:"email" binding:"" example:"zhangsan@example.com"` // 邮箱
}
type UserDataItem struct {
	ID        uint           `json:"id"`                                                            // ID
	CreatedAt string         `json:"createdAt"`                                                     // 创建时间
	UpdatedAt string         `json:"updatedAt"`                                                     // 更新时间
	Username  string         `json:"username" binding:"required" example:"zhangsan"`                // 用户名
	Nickname  string         `json:"nickname" binding:"required" example:"Jackal"`                  // 昵称
	Avatar    string         `json:"avatar" example:"https://example.com/avatar.jpg"`               // 头像
	Email     string         `json:"email" binding:"required,email" example:"zhangsan@example.com"` // 邮箱
	Phone     string         `json:"phone" binding:"" example:"+86-13966668888"`                    // 手机
	Status    int            `json:"status" example:"1"`                                            // 状态 0:待激活 1:正常 2:禁用
	Roles     []RoleDataItem `json:"roles"`                                                         // 角色
} // @name User
type UserSearchResponseData struct {
	List  []UserDataItem `json:"list"`  // 列表
	Total int64          `json:"total"` // 总数
} // @name UserList
type UserSearchResponse struct {
	Response
	Data UserSearchResponseData
}

// for Get
type UserResponse struct {
	Response
	Data UserDataItem
}

// for Create | Update
type UserRequest struct {
	Username string   `json:"username" binding:"required" example:"zhangsan"`
	Nickname string   `json:"nickname" binding:"" example:"Jackal"`
	Email    string   `json:"email" binding:"required" example:"zhangsan@example.com"`
	Phone    string   `json:"phone" binding:"" example:"+86-13966668888"`
	Status   int      `json:"status" binding:"" example:"1"` // 状态 0:待激活 1:正常 2:禁用
	Roles    []string `json:"roles" example:""`
}

// User Permission
type UserPermissionResponseData struct {
	List  []string `json:"list"`
	Total int64    `json:"total"`
}
type UserPermissionResponse struct {
	Response
	Data UserPermissionResponseData
}

// for Register
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"zhangsan@example.com"` // 邮箱
	Password string `json:"password" binding:"required" example:"123456"`                  // 密码
}

// for Login
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"zhangsan"`
	Password string `json:"password" binding:"required" example:"123456"`
}
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
}
type LoginResponse struct {
	Response
	Data LoginResponseData
}
