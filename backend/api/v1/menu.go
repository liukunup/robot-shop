package v1

// CRUD
type MenuSearchRequest struct {
	Page     int    `form:"page" binding:"required,min=1" example:"1"`              // 页码
	PageSize int    `form:"pageSize" binding:"required,min=1,max=100" example:"10"` // 分页大小
	Name     string `form:"name" example:"User"`                                    // 名称
	Path     string `form:"path" example:"/admin/user"`                             // 路径
	Access   string `form:"access" example:"canAdmin"`                              // 可见性
}
type MenuDataItem struct {
	ID                 uint   `json:"id,omitempty" example:"1"`                          // ID
	CreatedAt          string `json:"createdAt,omitempty" example:"2006-01-02 15:04:05"` // 创建时间
	UpdatedAt          string `json:"updatedAt,omitempty" example:"2006-01-02 15:04:05"` // 更新时间
	ParentID           uint   `json:"parentId,omitempty" example:"0"`                    // 父级菜单
	Icon               string `json:"icon,omitempty" example:"crown"`                    // 图标
	Name               string `json:"name,omitempty" example:"User"`                     // 名称
	Path               string `json:"path" example:"/admin/user"`                        // 路径
	Component          string `json:"component,omitempty" example:"@/pages/Admin/User"`  // 组件
	Access             string `json:"access,omitempty" example:"canAdmin"`               // 可见性
	Locale             string `json:"locale,omitempty"`                                  // 国际化
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
	List  []MenuDataItem `json:"list"`  // 列表
	Total int64          `json:"total"` // 总数
} // @name MenuList
type MenuSearchResponse struct {
	Response
	Data MenuSearchResponseData
}

type MenuResponse struct {
	Response
	Data MenuDataItem
}

type MenuRequest struct {
	ParentID           uint   `json:"parentId" example:"0"`                   // 父级菜单
	Icon               string `json:"icon" example:"crown"`                   // 图标
	Name               string `json:"name" example:"User"`                    // 名称
	Path               string `json:"path" example:"/admin/user"`             // 路径
	Component          string `json:"component" example:"@/pages/Admin/User"` // 组件
	Access             string `json:"access" example:"canAdmin"`              // 可见性
	Locale             string `json:"locale"`                                 // 国际化
	Redirect           string `json:"redirect"`                               // 重定向
	Target             string `json:"target"`                                 // 指定外链打开形式
	HideChildrenInMenu bool   `json:"hideChildrenInMenu"`                     // 隐藏子节点
	HideInMenu         bool   `json:"hideInMenu"`                             // 隐藏自身和子节点
	FlatMenu           bool   `json:"flatMenu"`                               // 隐藏自身+子节点提升并打平
	Disabled           bool   `json:"disabled"`
	Tooltip            string `json:"tooltip"`
	DisabledTooltip    bool   `json:"disabledTooltip"`
	Key                string `json:"key"`
	ParentKeys         string `json:"parentKeys"`
}

// Dynamic Menu
type MenuNode struct {
	MenuDataItem
	Children []*MenuNode `json:"children,omitempty"` // 子菜单
}
type DynamicMenuResponseData struct {
	List []*MenuNode `json:"list"` // 顶级菜单
}
type DynamicMenuResponse struct {
	Response
	Data DynamicMenuResponseData
}
