package v1

type PageRequest struct {
	Page    int                    `json:"page" form:"page" binding:"required,min=1"`
	Size    int                    `json:"size" form:"size" binding:"required,min=1,max=100"`
	Options map[string]interface{} `json:"options"`
}

type PageResponse[T any] struct {
	Total int64 `json:"total"`
	Items []T   `json:"items"`
}
