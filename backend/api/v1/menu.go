package v1

// Menu List
type MenuListRequest struct {
	Page     int `form:"page" binding:"required" example:"1"`      // 页码
	PageSize int `form:"pageSize" binding:"required" example:"10"` // 分页大小
}
type MenuDataItem struct {
	ID        uint   `json:"id,omitempty" example:"1"`                          // ID
	CreatedAt string `json:"createdAt,omitempty" example:"2006-01-02 15:04:05"` // 创建时间
	UpdatedAt string `json:"updatedAt,omitempty" example:"2006-01-02 15:04:05"` // 更新时间
	ParentID  uint   `json:"parentId,omitempty"`
	Path      string `json:"path"`
	Redirect  string `json:"redirect,omitempty"`
	Component string `json:"component,omitempty"`
	Name      string `json:"name,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Access    string `json:"access,omitempty"`
	Weight    int    `json:"weight,omitempty"`
} // @name Menu
type MenuListResponseData struct {
	List  []MenuDataItem `json:"list"`  // 列表
	Total int64          `json:"total"` // 总数
} // @name MenuList
type MenuListResponse struct {
	Response
	Data MenuListResponseData
}

// Menu Tree
type MenuDataNode struct {
	MenuDataItem
	Children []MenuDataItem `json:"children,omitempty"`
} // @name MenuTreeNode
type MenuTreeResponseData struct {
	Root []MenuDataNode `json:"root"`
}
type MenuTreeResponse struct {
	Response
	Data MenuTreeResponseData
}

// for Get
type MenuResponse struct {
	Response
	Data MenuDataItem
}

// for Create | Update
type MenuRequest struct {
	ParentID  uint   `json:"parentId,omitempty" example:"0"`
	Path      string `json:"path"`
	Redirect  string `json:"redirect,omitempty"`
	Component string `json:"component,omitempty"`
	Name      string `json:"name,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Access    string `json:"access,omitempty"`
	Weight    int    `json:"weight,omitempty" example:"0"`
}
