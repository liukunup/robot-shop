package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data"`
	Code     int         `json:"code,omitempty"`
	Message  string      `json:"message,omitempty"`
	ShowType int         `json:"showType,omitempty"`
}

const (
	SILENT       = 0 // 不提示
	WARNING      = 1 // 警告提示
	ERROR        = 2 // 错误提示
	NOTIFICATION = 3 // 通知提示
	REDIRECT     = 9 // 页面重定向
) // @name ShowType

func HandleSuccess(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = map[string]any{}
	}
	resp := Response{Success: true, Data: data, Code: errorCodeMap[ErrSuccess], Message: ErrSuccess.Error(), ShowType: SILENT}
	if _, ok := errorCodeMap[ErrSuccess]; !ok {
		resp = Response{Success: true, Data: data}
	}
	ctx.JSON(http.StatusOK, resp)
}

func HandleWithShowType(ctx *gin.Context, httpCode int, err error, data interface{}, showType int) {
	if data == nil {
		data = map[string]string{}
	}
	resp := Response{Success: false, Data: data, Code: errorCodeMap[err], Message: err.Error(), ShowType: showType}
	if _, ok := errorCodeMap[err]; !ok {
		resp = Response{Success: false, Data: data, Code: 500, Message: "unknown error", ShowType: ERROR}
	}
	ctx.JSON(httpCode, resp)
}

func HandleWarning(ctx *gin.Context, httpCode int, err error, data interface{}) {
	HandleWithShowType(ctx, httpCode, err, data, WARNING)
}

func HandleError(ctx *gin.Context, httpCode int, err error, data interface{}) {
	HandleWithShowType(ctx, httpCode, err, data, ERROR)
}

func HandleNotification(ctx *gin.Context, httpCode int, err error, data interface{}) {
	HandleWithShowType(ctx, httpCode, err, data, NOTIFICATION)
}

func HandleRedirect(ctx *gin.Context, httpCode int, err error, data interface{}) {
	HandleWithShowType(ctx, httpCode, err, data, REDIRECT)
}

type Error struct {
	Code    int
	Message string
}

var errorCodeMap = map[error]int{}

func newError(code int, msg string) error {
	err := errors.New(msg)
	errorCodeMap[err] = code
	return err
}

func (e Error) Error() string {
	return e.Message
}
