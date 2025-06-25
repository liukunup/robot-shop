package v1

// for Search
type RobotSearchRequest struct {
	Page     int    `form:"page" binding:"required" example:"1"`      // 页码
	PageSize int    `form:"pageSize" binding:"required" example:"10"` // 分页大小
	Name     string `form:"name" example:"bot"`                       // 筛选项: 名称 模糊匹配
	Desc     string `form:"desc" example:"robot"`                     // 筛选项: 描述 模糊匹配
	Owner    string `form:"owner" example:"Billy"`                    // 筛选项: 所有者 精确匹配
}
type RobotDataItem struct {
	Id        uint   `json:"id" example:"1"`                                  // ID
	CreatedAt string `json:"createdAt" example:"2006-01-02 15:04:05"`         // 创建时间
	UpdatedAt string `json:"updatedAt" example:"2006-01-02 15:04:05"`         // 更新时间
	Name      string `json:"name" example:"bot"`                              // 名称
	Desc      string `json:"desc" example:"It's a robot"`                     // 描述
	Webhook   string `json:"webhook" example:"https://example.com/webhook"`   // 通知地址
	Callback  string `json:"callback" example:"https://example.com/callback"` // 回调地址
	Enabled   bool   `json:"enabled" example:"true"`                          // 是否启用
	Owner     string `json:"owner" example:"Billy"`                           // 所有者
} // @name Robot
type RobotSearchResponseData struct {
	List  []RobotDataItem `json:"list"`  // 列表
	Total int64           `json:"total"` // 总数
} // @name RobotList
type RobotSearchResponse struct {
	Response
	Data RobotSearchResponseData
}

// for Get
type RobotResponse struct {
	Response
	Data RobotDataItem
}

// for Create | Update
type RobotRequest struct {
	Name     string `json:"name" example:"bot"`                              // 名称
	Desc     string `json:"desc" example:"It's a robot"`                     // 描述
	Webhook  string `json:"webhook" example:"https://example.com/webhook"`   // 通知地址
	Callback string `json:"callback" example:"https://example.com/callback"` // 回调地址
	Enabled  bool   `json:"enabled" example:"true"`                          // 是否启用
	Owner    string `json:"owner" example:"Billy"`                           // 所有者
}
