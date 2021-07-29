package e

var MsgFlags = map[int]string{
	// Success:                       "ok",
	// Error:                         "fail",
	// InvalidParams:                 "請求參數錯誤",
	// UsernameOrEmailAlreadyExist:   "使用者名稱或信箱已經存在",
	// UsernameAndPasswordCanNotNull: "使用者名稱及密碼不能為空",
	// EmailIsNotValid:               "請輸入正確的信箱格式",
	// UsernameNotFound:              "使用者名稱不存在",
	// PasswordDoesNotCorrect: "使用者密碼錯誤",
	// LastGameIsOver:                "前一局遊戲已經結束，請等待新的一局開始",
	// GamePlayerIsFull:              "此局遊戲玩家已滿，等待開獎",
	// UserNotFound:                  "使用者不存在",
	// BabyPayoutNotValid:            "寶貝號不存在(請輸入0-4)",
	// InsufficientBalance:           "遊戲餘額不足",
	// EmailOrPhoneCanNotNull:        "電子信箱或手機不能為空",
	// PhoneIsNotValid:               "請輸入正確的手機格式",
	// VerifyCodeIsInvalid:           "無效的驗證碼",
	// VerifyCodeIsTimeout:           "驗證碼已過期",
	// NotLogined:                    "請先登入",
	// LoginExpired:                  "請重新登入",
	// EmailOrPhoneAreOnLimit:        "註冊的手機或Email已達上限",
	// TokenInvalid:                  "無效的Token",
	// RecomandInvalid:               "無效的推薦碼",
	// WrongOldPassword:              "舊密碼錯誤",
	// DifferentPassword:             "兩次輸入的密碼不相同",
	// WaitForFiveMinute:             "請等候五分鐘再重試",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[Error]
}
