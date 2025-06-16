package v1

type GetRoleListRequest struct {
	Page     int    `form:"page" binding:"required" example:"1"`
	PageSize int    `form:"pageSize" binding:"required" example:"10"`
	Sid      string `form:"sid" binding:"" example:"1"`
	Name     string `form:"name" binding:"" example:"Admin"`
}
type RoleDataItem struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Sid       string `json:"sid"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}
type GetRolesResponseData struct {
	List  []RoleDataItem `json:"list"`
	Total int64          `json:"total"`
}
type GetRolesResponse struct {
	Response
	Data GetRolesResponseData
}
type RoleCreateRequest struct {
	Sid  string `form:"sid" binding:"required" example:"1"`
	Name string `form:"name" binding:"required" example:"Admin"`
}
type RoleUpdateRequest struct {
	ID   uint   `form:"id" binding:"required" example:"1"`
	Sid  string `form:"sid" binding:"required" example:"1"`
	Name string `form:"name" binding:"required" example:"Admin"`
}
type RoleDeleteRequest struct {
	ID uint `form:"id" binding:"required" example:"1"`
}