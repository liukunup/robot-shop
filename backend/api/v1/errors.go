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
	ErrEmailAlreadyUse    = newError(1001, "The email is already in use.")
	ErrUsernameAlreadyUse = newError(1002, "The username is already in use.")
)
