package v1

// CRUD
type ApiSearchRequest struct {
	Page     int    `form:"page" binding:"required,min=1" example:"1" `             // 页码
	PageSize int    `form:"pageSize" binding:"required,min=1,max=100" example:"10"` // 分页大小
	Group    string `form:"group" example:"User"`                                   // 筛选项: 分组 精确匹配
	Name     string `form:"name" example:"ListUsers"`                               // 筛选项: 名称 模糊匹配
	Path     string `form:"path" example:"/v1/admin/users"`                         // 筛选项: 路径 模糊匹配
	Method   string `form:"method" example:"GET"`                                   // 筛选项: 方法 精确匹配
}
type ApiDataItem struct {
	ID        uint   `json:"id,omitempty" example:"1"`                          // ID
	CreatedAt string `json:"createdAt,omitempty" example:"2006-01-02 15:04:05"` // 创建时间
	UpdatedAt string `json:"updatedAt,omitempty" example:"2006-01-02 15:04:05"` // 更新时间
	Group     string `json:"group" example:"User"`                              // 分组
	Name      string `json:"name" example:"ListUsers"`                          // 名称
	Path      string `json:"path" example:"/v1/admin/users"`                    // 路径
	Method    string `json:"method" example:"GET"`                              // 方法
} // @name Api
type ApiSearchResponseData struct {
	List  []ApiDataItem `json:"list"`  // 列表
	Total int64         `json:"total"` // 总数
} // @name ApiList
type ApiSearchResponse struct {
	Response
	Data ApiSearchResponseData
}

type ApiResponse struct {
	Response
	Data ApiDataItem
}

type ApiRequest struct {
	Group  string `json:"group" example:"User"`           // 分组
	Name   string `json:"name" example:"ListUsers"`       // 名称
	Path   string `json:"path" example:"/v1/admin/users"` // 路径
	Method string `json:"method" example:"GET"`           // 方法
}
