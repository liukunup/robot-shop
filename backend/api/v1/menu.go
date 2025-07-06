package v1

// for Search
type MenuSearchRequest struct {
	Page     int    `form:"page" binding:"required" example:"1"`      // 页码
	PageSize int    `form:"pageSize" binding:"required" example:"10"` // 分页大小
	Name     string `form:"name" example:"User"`                      // 名称
	Path     string `form:"path" example:"/admin/user"`               // 路径
	Access   string `form:"access" example:"canAdmin"`                // 可见性
}
type MenuDataItem struct {
	ID                 uint   `json:"id" example:"1"`                                    // ID
	CreatedAt          string `json:"createdAt,omitempty" example:"2006-01-02 15:04:05"` // 创建时间
	UpdatedAt          string `json:"updatedAt,omitempty" example:"2006-01-02 15:04:05"` // 更新时间
	ParentID           uint   `json:"parentId,omitempty" example:"0"`                    // 父级菜单
	Icon               string `json:"icon,omitempty" example:"crown"`                    // 图标
	Name               string `json:"name,omitempty" example:"User"`                     // 名称
	Path               string `json:"path" example:"/admin/user"`                        // 路径
	Component          string `json:"component,omitempty" example:"@/pages/Admin/User"`  // 组件
	Access             string `json:"access,omitempty" example:"canAdmin"`               // 可见性
	Locale             string `json:"locale,omitempty"`                                  // 本地化
	Weight             uint   `json:"weight,omitempty"`                                  // 权重
	Redirect           string `json:"redirect,omitempty"`                                // 重定向
	Target             string `json:"target,omitempty"`                                  // 指定外链打开形式
	HideChildrenInMenu bool   `json:"hideChildrenInMenu,omitempty"`                      // 隐藏子节点
	HideInMenu         bool   `json:"hideInMenu,omitempty"`                              // 隐藏自身和子节点
	FlatMenu           bool   `json:"flatMenu,omitempty"`                                // 隐藏自身+子节点提升并打平
	Disabled           bool   `json:"disabled,omitempty"`
	Tooltip            string `json:"tooltip,omitempty"`
	DisabledTooltip    bool   `json:"disabledTooltip,omitempty"`
	Key                string `json:"key,omitempty"`
	ParentKeys         string `json:"parentKeys,omitempty"`
} // @name Menu
type MenuSearchResponseData struct {
	List  []MenuDataItem `json:"list"`                         // 列表
	Total int64          `json:"total,omitempty" example:"10"` // 总数
} // @name MenuList
type MenuSearchResponse struct {
	Response
	Data MenuSearchResponseData
}

// for Dynamic Menu
type MenuNode struct {
	MenuDataItem
	Children []*MenuNode `json:"children,omitempty"` // 子节点 or 子菜单
}
type DynamicMenuResponseData struct {
	List []*MenuNode `json:"list"` // 顶级菜单
}
type DynamicMenuResponse struct {
	Response
	Data DynamicMenuResponseData
}

// for Get
type MenuResponse struct {
	Response
	Data MenuDataItem
}

// for Create | Update
type MenuRequest struct {
	ParentID           uint   `json:"parentId,omitempty"`           // 父级菜单
	Icon               string `json:"icon,omitempty"`               // 图标
	Name               string `json:"name,omitempty"`               // 名称
	Path               string `json:"path"`                         // 路径
	Component          string `json:"component,omitempty"`          // 组件
	Access             string `json:"access,omitempty"`             // 可见性
	Locale             string `json:"locale,omitempty"`             // 本地化
	Weight             uint   `json:"weight,omitempty"`             // 权重
	Redirect           string `json:"redirect,omitempty"`           // 重定向
	Target             string `json:"target,omitempty"`             // 指定外链打开形式
	HideChildrenInMenu bool   `json:"hideChildrenInMenu,omitempty"` // 隐藏子节点
	HideInMenu         bool   `json:"hideInMenu,omitempty"`         // 隐藏自身和子节点
	FlatMenu           bool   `json:"flatMenu,omitempty"`           // 隐藏自身+子节点提升并打平
	Disabled           bool   `json:"disabled,omitempty"`
	Tooltip            string `json:"tooltip,omitempty"`
	DisabledTooltip    bool   `json:"disabledTooltip,omitempty"`
	Key                string `json:"key,omitempty"`
	ParentKeys         string `json:"parentKeys,omitempty"`
}
