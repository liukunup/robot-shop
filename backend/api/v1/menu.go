package v1

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
}

type ListMenuResponseData struct {
	List []MenuDataItem `json:"list"`
}

type ListMenuResponse struct {
	Response
	Data ListMenuResponseData
}

type MenuCreateRequest struct {
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

type MenuUpdateRequest struct {
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
	UpdatedAt  string `json:"updatedAt"`
}

type MenuDeleteRequest struct {
	ID uint `form:"id"` // 唯一id，使用整数表示
}
