package e

const (
	Success       = 200
	InvalidParams = 400
	Unauthorized  = 401
	Forbidden     = 403
	Error         = 500
	SmsError      = "statuscode=e"
	DataNotExist  = iota + 10001
)
