package v1

// for Search
type UserSearchRequest struct {
	Page     int    `form:"page" binding:"required" example:"1"`      // 页码
	PageSize int    `form:"pageSize" binding:"required" example:"10"` // 分页大小
	Email    string `form:"email" example:"zhangsan@example.com"`     // 邮箱
	Username string `form:"username" example:"zhangsan"`              // 用户名
	Nickname string `form:"nickname" example:"Jackal"`                // 昵称
}
type UserDataItem struct {
	ID        uint           `json:"userid,omitempty" example:"1"`                              // ID
	CreatedAt string         `json:"createdAt,omitempty"  example:"2006-01-02 15:04:05"`        // 创建时间
	UpdatedAt string         `json:"updatedAt,omitempty"  example:"2006-01-02 15:04:05"`        // 更新时间
	Email     string         `json:"email,omitempty" example:"zhangsan@example.com"`            // 邮箱
	Username  string         `json:"username,omitempty" example:"zhangsan"`                     // 用户名
	Avatar    string         `json:"avatar,omitempty" example:"https://example.com/avatar.jpg"` // 头像
	Nickname  string         `json:"nickname,omitempty" example:"Jackal"`                       // 昵称
	Bio       string         `json:"bio,omitempty" example:"I am Jackal"`                       // 个人简介
	Language  string         `json:"language,omitempty" example:"zh-CN"`                        // 语言
	Timezone  string         `json:"timezone,omitempty" example:"Asia/Shanghai"`                // 时区
	Theme     string         `json:"theme,omitempty" example:"light"`                           // 主题
	Status    int            `json:"status,omitempty" example:"1"`                              // 状态 0:待激活 1:正常 2:禁用
	Roles     []RoleDataItem `json:"roles,omitempty"`                                           // 角色
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
	Email    string   `json:"email" example:"zhangsan@example.com"` // 邮箱
	Username string   `json:"username" example:"zhangsan"`          // 用户名
	Nickname string   `json:"nickname" example:"Jackal"`            // 昵称
	Bio      string   `json:"bio" example:"I am Jackal"`            // 个人简介
	Language string   `json:"language" example:"zh-CN"`             // 语言
	Timezone string   `json:"timezone" example:"Asia/Shanghai"`     // 时区
	Theme    string   `json:"theme" example:"light"`                // 主题
	Status   int      `json:"status" example:"1"`                   // 状态 0:待激活 1:正常 2:禁用
	Roles    []string `json:"roles"`                                // 角色
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
	Username string `json:"username" binding:"required" example:"zhangsan"` // 用户名
	Password string `json:"password" binding:"required" example:"123456"`   // 密码
}
type LoginResponseData struct {
	AccessToken string `json:"accessToken"` // 访问令牌
}
type LoginResponse struct {
	Response
	Data LoginResponseData
}

// for UpdatePassword
type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required" example:"123456"` // 旧密码
	NewPassword string `json:"newPassword" binding:"required" example:"123456"` // 新密码
}

// for ResetPassword
type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"zhangsan@example.com"` // 邮箱
}
