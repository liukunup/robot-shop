package v1

type RobotSearchRequest struct {
	Page     int    `form:"page" binding:"required" example:"1"`      // 页码
	PageSize int    `form:"pageSize" binding:"required" example:"10"` // 分页大小
	Name     string `form:"name" example:"bot"`                       // 筛选项: 名称 模糊匹配
} // @name RobotSearchParams
type RobotData struct {
	Id        uint   `json:"id" example:"1"`                          // ID
	CreatedAt string `json:"createdAt" example:"2006-01-02 15:04:05"` // 创建时间
	UpdatedAt string `json:"updatedAt" example:"2006-01-02 15:04:05"` // 更新时间
	Name      string `json:"name" example:"bot"`
	Desc      string `json:"desc" example:"It's a robot"`
	Webhook   string `json:"webhook" example:"https://example.com/webhook"`
	Callback  string `json:"callback" example:"https://example.com/callback"`
	Enabled   bool   `json:"enabled" example:"true"`
	Owner     string `json:"owner" example:"Billy"`
} // @name Robot
type RobotSearchResponseData struct {
	List  []RobotData `json:"list"`
	Total int64       `json:"total"`
} // @name RobotList
type RobotSearchResponse struct {
	Response
	Data RobotSearchResponseData
}

type RobotResponse struct {
	Response
	Data RobotData
}

type RobotRequest struct {
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Webhook  string `json:"webhook"`
	Callback string `json:"callback"`
	Enabled  bool   `json:"enabled"`
	Owner    string `json:"owner"`
} // @name RobotParams
