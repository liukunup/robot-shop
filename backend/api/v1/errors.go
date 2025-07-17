package v1

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
