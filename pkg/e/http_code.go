package e

const (
	Success        = 200
	InvalidParams  = 400
	Unauthorized   = 401
	Forbidden      = 403
	StatusNotFound = 404
	Error          = 500

	// auth
	ResetRedirect   = 10001
	AccountNotExist = 90001
	TokenNotExist   = 90002
	TokenExpired    = 90003
	PasswordNoMatch = 90004
)
