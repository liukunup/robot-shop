package v1

// for Search
type UserSearchRequest struct {
	Page     int    `form:"page" binding:"required" example:"1"`      // 页码
	PageSize int    `form:"pageSize" binding:"required" example:"10"` // 分页大小
	Username string `form:"username" example:"zhangsan"`              // 用户名
	Nickname string `form:"nickname" example:"Jackal"`                // 昵称
	Phone    string `form:"phone" example:"13966668888"`              // 手机
	Email    string `form:"email" example:"zhangsan@example.com"`     // 邮箱
}
type UserDataItem struct {
	ID        uint           `json:"id"`                                                        // ID
	CreatedAt string         `json:"createdAt,omitempty"  example:"2006-01-02 15:04:05"`        // 创建时间
	UpdatedAt string         `json:"updatedAt,omitempty"  example:"2006-01-02 15:04:05"`        // 更新时间
	Username  string         `json:"username" example:"zhangsan"`                               // 用户名
	Nickname  string         `json:"nickname,omitempty" example:"Jackal"`                       // 昵称
	Avatar    string         `json:"avatar,omitempty" example:"https://example.com/avatar.jpg"` // 头像
	Email     string         `json:"email" example:"zhangsan@example.com"`                      // 邮箱
	Phone     string         `json:"phone,omitempty" example:"13966668888"`                     // 手机
	Status    int            `json:"status" example:"1"`                                        // 状态 0:待激活 1:正常 2:禁用
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
	Username string   `json:"username" binding:"required" example:"zhangsan"`
	Nickname string   `json:"nickname" example:"Jackal"`
	Email    string   `json:"email" binding:"required" example:"zhangsan@example.com"`
	Phone    string   `json:"phone" example:"13966668888"`
	Status   int      `json:"status" example:"1"` // 状态 0:待激活 1:正常 2:禁用
	Roles    []string `json:"roles"`
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

// for UpdatePassword
type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required" example:"123456"` // 旧密码
	NewPassword string `json:"newPassword" binding:"required" example:"123456"` // 新密码
}

// for ResetPassword
type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"zhangsan@example.com"` // 邮箱
}
