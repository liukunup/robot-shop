package v1

// for Search
type MenuSearchRequest struct {
	Page     int `form:"page" binding:"required" example:"1"`      // 页码
	PageSize int `form:"pageSize" binding:"required" example:"10"` // 分页大小
}
type MenuDataItem struct {
	ID         uint   `json:"id,omitempty"`         // 唯一id，使用整数表示
	ParentID   uint   `json:"parentId,omitempty"`   // 父级菜单的id，使用整数表示
	Weight     int    `json:"weight"`               // 排序权重
	Path       string `json:"path"`                 // 地址
	Title      string `json:"title"`                // 展示名称
	Name       string `json:"name,omitempty"`       // 同路由中的name，唯一标识
	Component  string `json:"component,omitempty"`  // 绑定的组件
	Locale     string `json:"locale,omitempty"`     // 本地化标识
	Icon       string `json:"icon,omitempty"`       // 图标，使用字符串表示
	Redirect   string `json:"redirect,omitempty"`   // 重定向地址
	KeepAlive  bool   `json:"keepAlive,omitempty"`  // 是否保活
	HideInMenu bool   `json:"hideInMenu,omitempty"` // 是否保活
	URL        string `json:"url,omitempty"`        // iframe模式下的跳转url，不能与path重复
	UpdatedAt  string `json:"updatedAt,omitempty"`  // 是否保活
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
	ParentID   uint   `json:"parentId,omitempty"`   // 父级菜单的id，使用整数表示
	Weight     int    `json:"weight"`               // 排序权重
	Path       string `json:"path"`                 // 地址
	Title      string `json:"title"`                // 展示名称
	Name       string `json:"name,omitempty"`       // 同路由中的name，唯一标识
	Component  string `json:"component,omitempty"`  // 绑定的组件
	Locale     string `json:"locale,omitempty"`     // 本地化标识
	Icon       string `json:"icon,omitempty"`       // 图标，使用字符串表示
	Redirect   string `json:"redirect,omitempty"`   // 重定向地址
	KeepAlive  bool   `json:"keepAlive,omitempty"`  // 是否保活
	HideInMenu bool   `json:"hideInMenu,omitempty"` // 是否保活
	URL        string `json:"url,omitempty"`        // iframe模式下的跳转url，不能与path重复
}
