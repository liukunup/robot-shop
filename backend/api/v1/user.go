package v1

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"zhangsan@example.com"`
	Password string `json:"password" binding:"required" example:"123456"`
} // @name RegisterParams

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"zhangsan"`
	Password string `json:"password" binding:"required" example:"123456"`
} // @name LoginParams

type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
} // @name LoginResult

type LoginResponse struct {
	Response
	Data LoginResponseData
}

type ListUsersRequest struct {
	Page     int    `form:"page" binding:"required" example:"1"`
	PageSize int    `form:"pageSize" binding:"required" example:"10"`
	Username string `json:"username" binding:"" example:"zhangsan"`
	Nickname string `json:"nickname" binding:"" example:"法外狂徒"`
	Phone    string `form:"phone" binding:"" example:"+86-13966668888"`
	Email    string `form:"email" binding:"" example:"zhangsan@example.com"`
}

type UserDataItem struct {
	UserID    uint     `json:"id"`
	Username  string   `json:"username" binding:"required" example:"zhangsan"`
	Nickname  string   `json:"nickname" binding:"required" example:"法外狂徒"`
	Email     string   `json:"email" binding:"required,email" example:"zhangsan@example.com"`
	Phone     string   `form:"phone" binding:"" example:"+86-13966668888"`
	Roles     []string `json:"roles" example:""`
	UpdatedAt string   `json:"updatedAt"`
	CreatedAt string   `json:"createdAt"`
}

type ListUsersResponseData struct {
	List  []UserDataItem `json:"list"`
	Total int64          `json:"total"`
}

type ListUsersResponse struct {
	Response
	Data ListUsersResponseData
}

type UserCreateRequest struct {
	Username string   `json:"username" binding:"required" example:"zhangsan"`
	Nickname string   `json:"nickname" binding:"" example:"法外狂徒"`
	Password string   `json:"password" binding:"required" example:"123456"`
	Email    string   `json:"email" binding:"" example:"zhangsan@example.com"`
	Phone    string   `form:"phone" binding:"" example:"+86-13966668888"`
	Roles    []string `json:"roles" example:""`
}

type UserUpdateRequest struct {
	UserID   uint     `json:"id"`
	Username string   `json:"username" binding:"required" example:"zhangsan"`
	Nickname string   `json:"nickname" binding:"" example:"法外狂徒"`
	Password string   `json:"password" binding:"" example:"123456"`
	Email    string   `json:"email" binding:"" example:"zhangsan@example.com"`
	Phone    string   `form:"phone" binding:"" example:"+86-13966668888"`
	Roles    []string `json:"roles" example:""`
}

type UserDeleteRequest struct {
	UserID uint `form:"id" binding:"required" example:"1"`
}

type GetUserResponseData struct {
	UserID    uint     `json:"userid"`
	Username  string   `json:"username" example:"zhangsan"`
	Nickname  string   `json:"nickname" example:"法外狂徒"`
	Avatar    string   `json:"avatar" example:"https://example.com/avatar.jpg"`
	Email     string   `json:"email,omitempty" example:"zhangsan@example.com"`
	Phone     string   `json:"phone,omitempty" example:"+86-13966668888"`
	Roles     []string `json:"roles" example:""`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
} // @name CurrentUser

type GetUserResponse struct {
	Response
	Data GetUserResponseData
}
