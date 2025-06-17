package v1

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"zhangsan@163.com"`
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
	Username string `json:"username" binding:"" example:"张三"`
	Nickname string `json:"nickname" binding:"" example:"小Baby"`
	Phone    string `form:"phone" binding:"" example:"1858888888"`
	Email    string `form:"email" binding:"" example:"1234@gmail.com"`
}

type UserDataItem struct {
	ID        uint     `json:"id"`
	Username  string   `json:"username" binding:"required" example:"张三"`
	Nickname  string   `json:"nickname" binding:"required" example:"小Baby"`
	Email     string   `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Phone     string   `form:"phone" binding:"" example:"1858888888"`
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
	Username string   `json:"username" binding:"required" example:"张三"`
	Nickname string   `json:"nickname" binding:"" example:"小Baby"`
	Password string   `json:"password" binding:"required" example:"123456"`
	Email    string   `json:"email" binding:"" example:"1234@gmail.com"`
	Phone    string   `form:"phone" binding:"" example:"1858888888"`
	Roles    []string `json:"roles" example:""`
}

type UserUpdateRequest struct {
	ID       uint     `json:"id"`
	Username string   `json:"username" binding:"required" example:"张三"`
	Nickname string   `json:"nickname" binding:"" example:"小Baby"`
	Password string   `json:"password" binding:"" example:"123456"`
	Email    string   `json:"email" binding:"" example:"1234@gmail.com"`
	Phone    string   `form:"phone" binding:"" example:"1858888888"`
	Roles    []string `json:"roles" example:""`
}

type UserDeleteRequest struct {
	ID uint `form:"id" binding:"required" example:"1"`
}

type GetUserResponseData struct {
	ID        uint     `json:"id"`
	Username  string   `json:"username" example:"张三"`
	Nickname  string   `json:"nickname" example:"小Baby"`
	Email     string   `json:"email" example:"1234@gmail.com"`
	Phone     string   `form:"phone" example:"1858888888"`
	Avatar    string   `json:"avatar" example:"https://example.com/avatar.jpg"`
	Roles     []string `json:"roles" example:""`
	UpdatedAt string   `json:"updatedAt"`
	CreatedAt string   `json:"createdAt"`
} // @name CurrentUser

type GetUserResponse struct {
	Response
	Data GetUserResponseData
}
