package v1

// CRUD
type RobotSearchRequest struct {
	Page     int    `form:"page" binding:"required,min=1" example:"1"`               // 页码
	PageSize int    `form:"pageSize" binding:"required,min=1,max=1000" example:"10"` // 分页大小
	Name     string `form:"name" example:"bot"`                                      // 筛选项: 名称 模糊匹配
	Desc     string `form:"desc" example:"robot"`                                    // 筛选项: 描述 模糊匹配
	Owner    string `form:"owner" example:"Zhangsan"`                                // 筛选项: 所有者 精确匹配
}
type RobotDataItem struct {
	Id        uint   `json:"id,omitempty" example:"1"`                                  // ID
	CreatedAt string `json:"createdAt,omitempty" example:"2006-01-02 15:04:05"`         // 创建时间
	UpdatedAt string `json:"updatedAt,omitempty" example:"2006-01-02 15:04:05"`         // 更新时间
	Name      string `json:"name" example:"robot"`                                      // 名称
	Desc      string `json:"desc,omitempty" example:"it's a chatbot"`                   // 描述
	Webhook   string `json:"webhook,omitempty" example:"https://example.com/webhook"`   // 通知地址
	Callback  string `json:"callback,omitempty" example:"https://example.com/callback"` // 回调地址
	Enabled   bool   `json:"enabled" example:"true"`                                    // 是否启用
	Owner     string `json:"owner,omitempty" example:"Zhangsan"`                        // 所有者
} // @name Robot
type RobotSearchResponseData struct {
	List  []RobotDataItem `json:"list"`  // 列表
	Total int64           `json:"total"` // 总数
} // @name RobotList
type RobotSearchResponse struct {
	Response
	Data RobotSearchResponseData
}

type RobotResponse struct {
	Response
	Data RobotDataItem
}

type RobotRequest struct {
	Name     string `json:"name" example:"robot"`                            // 名称
	Desc     string `json:"desc" example:"it's a chatbot"`                   // 描述
	Webhook  string `json:"webhook" example:"https://example.com/webhook"`   // 通知地址
	Callback string `json:"callback" example:"https://example.com/callback"` // 回调地址
	Enabled  bool   `json:"enabled" example:"true"`                          // 是否启用
	Owner    string `json:"owner" example:"Zhangsan"`                        // 所有者
}
