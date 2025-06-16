package v1

type GetUserPermissionsData struct {
	List []string `json:"list"`
}

type GetRolePermissionsData struct {
	List []string `json:"list"`
}

type GetRolePermissionsRequest struct {
	Role string `form:"role" binding:"required" example:"admin"`
}

type UpdateRolePermissionRequest struct {
	Role string   `form:"role" binding:"required" example:"admin"`
	List []string `form:"list" binding:"required" example:""`
}
