package v1

// List
type GetRobotListRequest struct {
	Page     int    `form:"current" binding:"required" example:"1"`   // 页码
	PageSize int    `form:"pageSize" binding:"required" example:"10"` // 分页大小
	Name     string `form:"name" example:"bot"`                       // 筛选项: 名称 模糊匹配
} // @name GetRobotListParams
type RobotDataItem struct {
	Id        uint   `json:"id" example:"1"`                          // ID
	CreatedAt string `json:"createdAt" example:"2006-01-02 15:04:05"` // 创建时间
	UpdatedAt string `json:"updatedAt" example:"2006-01-02 15:04:05"` // 更新时间
	Name      string `json:"name" example:"bot"`
	Desc      string `json:"desc" example:"It's a robot"`
	Webhook   string `json:"webhook" example:"https://webhook.example.com"`
	Callback  string `json:"callback" example:"https://callback.example.com"`
	Enabled   bool   `json:"enabled" example:"true"`
	Owner     string `json:"owner" example:"Billy"`
} // @name Robot
type GetRobotListResponseData struct {
	List  []RobotDataItem `json:"list"`
	Total int64           `json:"total"`
} // @name RobotList
type GetRobotListResponse struct {
	Response
	Data GetRobotListResponseData
}

// Get
type GetRobotRequest struct {
	ID uint `json:"id" binding:"required" example:"1"` // ID
} // @name GetRobotParams
type GetRobotResponse struct {
	Response
	Data RobotDataItem
}

// Create
type RobotCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Desc     string `json:"desc"`
	Webhook  string `json:"webhook"`
	Callback string `json:"callback"`
	Enabled  bool   `json:"enabled"`
	Owner    string `json:"owner"`
} // @name RobotCreateParams

// Update
type RobotUpdateRequest struct {
	ID       uint   `json:"id" binding:"required" example:"1"` // ID
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Webhook  string `json:"webhook"`
	Callback string `json:"callback"`
	Enabled  bool   `json:"enabled"`
	Owner    string `json:"owner"`
} // @name RobotUpdateParams

// Delete
type RobotDeleteRequest struct {
	ID uint `json:"id" binding:"required" example:"1"` // ID
} // @name RobotDeleteParams
