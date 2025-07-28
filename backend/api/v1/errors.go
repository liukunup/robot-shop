package v1

/*
 * <错误码> 规则如下：
 * 1. 约定俗成，如 0 表示成功，其他表示失败；
 * 2. 分段管理，如 HTTP 状态码、业务错误码、系统错误码；
 * 3. 易理解、可自行解决，如 401 表示未授权，403 表示禁止访问，404 表示未找到，用户可以通过账号登录来解决等；
 */

var (
	// common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrForbidden           = newError(403, "Forbidden")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")
	ErrServiceUnavailable  = newError(503, "Service Unavailable")

	// more biz errors
	ErrEmailAlreadyUse         = newError(1001, "The email is already in use.")
	ErrUsernameAlreadyUse      = newError(1002, "The username is already in use.")
	ErrEmptyToken              = newError(1003, "token is empty")
	ErrInvalidToken            = newError(1004, "invalid token")
	ErrUnexpectedClaim         = newError(1005, "unexpected claims type")
	ErrTokenExpired            = newError(1006, "token has expired")
	ErrInvalidSigningMethod    = newError(1007, "invalid signing method")
	ErrInvalidKeyLength        = newError(1008, "invalid key length")
	ErrRedisUnavailable        = newError(1009, "redis service unavailable")
	ErrUnexpectedSigningMethod = newError(1010, "unexpected signing method")
	ErrInvalidAccessToken      = newError(1011, "invalid access token")
	ErrInvalidRefreshToken     = newError(1012, "invalid refresh token")
	ErrTokenAlreadyRevoked     = newError(1013, "token already revoked with later expiry")
	ErrAvatarSizeExceeded      = newError(1014, "avatar size exceeded")
	ErrAvatarTypeInvalid       = newError(1015, "avatar type invalid")
)
