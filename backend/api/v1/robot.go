package v1

type RobotResponseData struct {
	Id        uint                   `json:"id" example:"1"`
	RobotId   string                 `json:"robot_id" example:"454e7080-d105-4515-9c25-3e6fd9df176c"`
	Name      string                 `json:"name" example:"bot"`
	Desc      string                 `json:"desc" example:"it's a robot"`
	Webhook   string                 `json:"webhook" example:"https://webhook.example.com"`
	Callback  string                 `json:"callback" example:"https://callback.example.com"`
	Options   map[string]interface{} `json:"options" example:"{\"key\": \"value\"}"`
	Enabled   bool                   `json:"enabled" example:"true"`
	Owner     string                 `json:"owner" example:"Billy"`
	CreatedAt string                 `json:"createdAt" example:"2006-01-02 15:04:05"`
	UpdatedAt string                 `json:"updatedAt" example:"2006-01-02 15:04:05"`
}

type RobotRequest struct {
	Name     string                 `json:"name" binding:"required"`
	Desc     string                 `json:"desc"`
	Webhook  string                 `json:"webhook"`
	Callback string                 `json:"callback"`
	Options  map[string]interface{} `json:"options"`
	Enabled  bool                   `json:"enabled"`
	Owner    string                 `json:"owner"`
}
