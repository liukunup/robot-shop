package v1

// for Search
type MenuSearchRequest struct {
	Page     int `form:"page" binding:"required" example:"1"`      // 页码
	PageSize int `form:"pageSize" binding:"required" example:"10"` // 分页大小
}
type MenuDataItem struct {
	ID        uint   `json:"id,omitempty" example:"1"`                          // ID
	CreatedAt string `json:"createdAt,omitempty" example:"2006-01-02 15:04:05"` // 创建时间
	UpdatedAt string `json:"updatedAt,omitempty" example:"2006-01-02 15:04:05"` // 更新时间
	ParentID  uint   `json:"parentId,omitempty" example:"0"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Name      string `json:"name,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Weight    int    `json:"weight,omitempty" example:"0"`
} // @name Menu
type MenuSearchResponseData struct {
	List  []MenuDataItem `json:"list"`  // 列表
	Total int64          `json:"total"` // 总数
} // @name MenuList
type MenuSearchResponse struct {
	Response
	Data MenuSearchResponseData
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
	Component string `json:"component"`
	Name      string `json:"name,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Weight    int    `json:"weight,omitempty" example:"0"`
}
