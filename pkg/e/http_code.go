package e

const (
	Success                   = 200
	InvalidParams             = 400
	Unauthorized              = 401
	Forbidden                 = 403
	StatusNotFound            = 404
	StatusInternalServerError = 500

	// auth
	ResetRedirect   = 10001
	AccountNotExist = 90001
	TokenNotExist   = 90002
	TokenExpired    = 90003
	PasswordNoMatch = 90004

	// auth frontend
	MemberNotExist         = 91000
	EmailDuplicate         = 91001
	NoLogginOrTokenExpired = 91002

	// order frontend
	ProductPriceChange = 91100
)
