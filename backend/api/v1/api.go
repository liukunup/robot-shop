package v1

type ListApisRequest struct {
	Page     int    `form:"page" binding:"required" example:"1"`
	PageSize int    `form:"pageSize" binding:"required" example:"10"`
	Group    string `form:"group" binding:"" example:"权限管理"`
	Name     string `form:"name" binding:"" example:"菜单列表"`
	Path     string `form:"path" binding:"" example:"/v1/test"`
	Method   string `form:"method" binding:"" example:"GET"`
}

type ApiDataItem struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Method    string `json:"method"`
	Group     string `json:"group"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

type ListApisResponseData struct {
	List   []ApiDataItem `json:"list"`
	Total  int64         `json:"total"`
	Groups []string      `json:"groups"`
}

type ListApisResponse struct {
	Response
	Data ListApisResponseData
}

type ApiCreateRequest struct {
	Group  string `form:"group" binding:"" example:"权限管理"`
	Name   string `form:"name" binding:"" example:"菜单列表"`
	Path   string `form:"path" binding:"" example:"/v1/test"`
	Method string `form:"method" binding:"" example:"GET"`
}

type ApiUpdateRequest struct {
	ID     uint   `form:"id" binding:"required" example:"1"`
	Group  string `form:"group" binding:"" example:"权限管理"`
	Name   string `form:"name" binding:"" example:"菜单列表"`
	Path   string `form:"path" binding:"" example:"/v1/test"`
	Method string `form:"method" binding:"" example:"GET"`
}

type ApiDeleteRequest struct {
	ID uint `form:"id" binding:"required" example:"1"`
}
