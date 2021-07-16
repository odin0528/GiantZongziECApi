package e

const (
	Success                     = 200
	Error                       = 500
	InvalidParams               = 400
	SmsError                    = "statuscode=e"
	UsernameOrEmailAlreadyExist = iota + 10001
	UsernameAndPasswordCanNotNull
	EmailIsNotValid
	UsernameNotFound
	PasswordDoesNotCorrect
	LastGameIsOver
	GamePlayerIsFull
	UserNotFound
	BabyPayoutNotValid
	InsufficientBalance
	EmailOrPhoneCanNotNull
	PhoneIsNotValid
	VerifyCodeIsInvalid
	VerifyCodeIsTimeout
	NotLogined
	LoginExpired
	EmailOrPhoneAreOnLimit
	TokenInvalid
	RecomandInvalid
	WrongOldPassword
	DifferentPassword
	WaitForFiveMinute
)
