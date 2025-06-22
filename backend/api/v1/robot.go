package v1

// List
type ListRobotRequest struct {
	Page     int    `form:"page" binding:"required" example:"1"`
	PageSize int    `form:"pageSize" binding:"required" example:"10"`
	Name     string `form:"name" example:"bot"`
} // @name ListRobotParams
type RobotDataItem struct {
	Id        uint   `json:"id" example:"1"`
	Name      string `json:"name" example:"bot"`
	Desc      string `json:"desc" example:"It's a robot"`
	Webhook   string `json:"webhook" example:"https://webhook.example.com"`
	Callback  string `json:"callback" example:"https://callback.example.com"`
	Enabled   bool   `json:"enabled" example:"true"`
	Owner     string `json:"owner" example:"Billy"`
	CreatedAt string `json:"createdAt" example:"2006-01-02 15:04:05"`
	UpdatedAt string `json:"updatedAt" example:"2006-01-02 15:04:05"`
} // @name Robot
type ListRobotResponseData struct {
	List  []RobotDataItem `json:"list"`
	Total int64           `json:"total"`
} // @name RobotList
type ListRobotResponse struct {
	Response
	Data ListRobotResponseData
}

// Get
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
} // @name CreateRobotParams

// Update
type RobotUpdateRequest struct {
	ID       uint   `form:"id" binding:"required" example:"1"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Webhook  string `json:"webhook"`
	Callback string `json:"callback"`
	Enabled  bool   `json:"enabled"`
	Owner    string `json:"owner"`
} // @name UpdateRobotParams
