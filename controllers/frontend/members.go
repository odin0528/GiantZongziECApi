package frontend

import (
	"eCommerce/internal/auth"
	"eCommerce/pkg/e"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	. "eCommerce/internal/database"
	models "eCommerce/models/frontend"

	fb "github.com/huandu/facebook/v2"
	"github.com/liudng/godump"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/scrypt"
)

func GetMemberOrders(c *gin.Context) {
	g := Gin{c}
	var req *models.OrderQuery
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	MemeberID, _ := c.Get("member_id")
	req.MemberID = MemeberID.(int)

	orders, pagination := req.FetchAll()

	for index := range orders {
		orders[index].GetProducts()
	}

	g.PaginationResponse(http.StatusOK, e.Success, orders, pagination)
}

func GetMemberDelivery(c *gin.Context) {
	g := Gin{c}
	MemberID, _ := c.Get("member_id")
	query := &models.MemberDeliveryQuery{
		MemberID: MemberID.(int),
	}
	deliveries := query.FetchAll()

	g.Response(http.StatusOK, e.Success, deliveries)
}

func MemberFetch(c *gin.Context) {
	g := Gin{c}
	if c.Request.Header.Get("Authorization") == "" {
		g.Response(http.StatusOK, e.NoLogginOrTokenExpired, nil)
		return
	}

	tokenClaims, err := jwt.ParseWithClaims(c.Request.Header.Get("Authorization"), &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SIGN")), nil
	})

	if err != nil {
		g.Response(http.StatusOK, e.NoLogginOrTokenExpired, nil)
		return
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*models.Claims); ok && tokenClaims.Valid {
			g.Response(http.StatusOK, e.Success, models.Members{Nickname: claims.Nickname, Avatar: claims.Avatar})
		} else {
			g.Response(http.StatusOK, e.NoLogginOrTokenExpired, nil)
			return
		}
	}
}

func MemberLogin(c *gin.Context) {
	g := Gin{c}
	var req models.LoginReq
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	PlatformID, _ := c.Get("platform_id")

	query := models.MemberQuery{
		PlatformID: PlatformID.(int),
		Email:      req.Email,
	}

	member := query.Fetch()

	if member.ID == 0 {
		g.Response(http.StatusOK, e.MemberNotExist, err)
		return
	}

	dk, err := scrypt.Key([]byte(req.Password), auth.Salt, 1<<15, 8, 1, 64)
	if base64.StdEncoding.EncodeToString(dk) != member.Password {
		g.Response(http.StatusOK, e.MemberNotExist, nil)
		return
	}

	token := models.GenerateToken(member)

	g.Response(http.StatusOK, e.Success, map[string]interface{}{"token": token, "member": member})
}

func MemberOAuth(c *gin.Context) {
	g := Gin{c}

	var req models.OAuthReq
	var member models.Members

	PlatformID, _ := c.Get("platform_id")

	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	if req.Platform == "fb" {
		var user models.FbUser
		res, _ := fb.Get("/me?fields=id,name,gender,email,birthday,picture.type(large)", fb.Params{
			"access_token": req.Token,
		})

		res.Decode(&user)

		query := models.MemberQuery{
			OAuthUserID:   user.ID,
			OAuthPlatform: req.Platform,
			PlatformID:    PlatformID.(int),
		}

		member = query.Fetch()

		if member.ID == 0 {
			member.Nickname = user.Name
			member.Email = user.Email
			birthday, _ := time.Parse("01/02/2006", user.Birthday)
			member.Birthday = birthday
			member.OAuthPlatform = req.Platform
			member.OAuthUserID = user.ID
			member.PlatformID = PlatformID.(int)
			member.Avatar = user.Picture.Data.Url

			if user.Gender == "male" {
				member.Gender = 1
			} else if user.Gender == "female" {
				member.Gender = 0
			}

			DB.Create(&member)
		}
	} else if req.Platform == "line" {

		var result map[string]interface{}
		var userData map[string]interface{}

		params := url.Values{}
		params.Add("grant_type", "authorization_code")
		params.Add("code", req.Token)
		params.Add("redirect_uri", "https://example.com:3000/oauth/line")
		params.Add("client_id", "1656498603")
		params.Add("client_secret", "dc9e16a0b2a0d05c6a9870bc60ab7f9c")
		body := strings.NewReader(params.Encode())

		curl, _ := http.NewRequest("POST", "https://api.line.me/oauth2/v2.1/token", body)
		curl.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, _ := http.DefaultClient.Do(curl)
		rbody, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(rbody, &result)

		params = url.Values{}
		params.Add("id_token", result["id_token"].(string))
		params.Add("client_id", "1656498603")
		body = strings.NewReader(params.Encode())
		curl, _ = http.NewRequest("POST", "https://api.line.me/oauth2/v2.1/verify", body)
		curl.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, _ = http.DefaultClient.Do(curl)
		rbody, _ = ioutil.ReadAll(resp.Body)
		json.Unmarshal(rbody, &userData)

		godump.Dump(userData)

		defer resp.Body.Close()

		query := models.MemberQuery{
			OAuthUserID:   userData["sub"].(string),
			OAuthPlatform: req.Platform,
			PlatformID:    PlatformID.(int),
		}

		member = query.Fetch()

		if member.ID == 0 {
			if _, ok := userData["email"]; ok {
				member.Email = userData["email"].(string)
			}
			member.Nickname = userData["name"].(string)
			member.OAuthPlatform = req.Platform
			member.OAuthUserID = userData["sub"].(string)
			member.PlatformID = PlatformID.(int)
			member.Avatar = userData["picture"].(string)
			DB.Create(&member)
		}
	}

	token := models.GenerateToken(member)

	g.Response(http.StatusOK, e.Success, map[string]interface{}{"token": token, "member": member})
}

func MemberRegister(c *gin.Context) {
	g := Gin{c}
	var req models.MemberRegisterReq
	var birthday time.Time
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusOK, e.StatusInternalServerError, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		g.Response(http.StatusOK, e.InvalidParams, err.(validator.ValidationErrors)[0].Field())
		return
	}

	if req.Birthday != "" {
		birthday, err = time.Parse("20060102", req.Birthday)
		if err != nil {
			g.Response(http.StatusOK, e.InvalidParams, "Birthday")
			return
		}
	}

	PlatformID, _ := c.Get("platform_id")
	dk, err := scrypt.Key([]byte(req.Password), auth.Salt, 1<<15, 8, 1, 64)

	member := models.Members{
		PlatformID: PlatformID.(int),
		Email:      req.Email,
		Nickname:   req.Nickname,
		Phone:      req.Phone,
		Birthday:   birthday,
		Password:   base64.StdEncoding.EncodeToString(dk),
	}

	err = DB.Create(&member).Error

	if err != nil {
		g.Response(http.StatusOK, e.EmailDuplicate, err)
		return
	}

	/* token := uuid.New().String()
	memberToken := models.MemberToken{
		MemberID:   member.ID,
		Token:      token,
		PlatformID: PlatformID.(int),
	}
	DB.Create(&memberToken) */

	token := models.GenerateToken(member)

	// req.Create()
	g.Response(http.StatusOK, e.Success, map[string]interface{}{"token": token, "member": member})

}

func MemberLogout(c *gin.Context) {
	g := Gin{c}
	MemberID, _ := c.Get("member_id")
	memberToken := models.MemberToken{
		MemberID: MemberID.(int),
	}
	memberToken.CancelOldToken()

	g.Response(http.StatusOK, e.Success, nil)
}

func MemberDeliveryModify(c *gin.Context) {
	g := Gin{c}
	var delivery *models.MemberDelivery
	err := c.BindJSON(&delivery)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	PlatformID, _ := c.Get("platform_id")
	MemberID, _ := c.Get("member_id")
	delivery.PlatformID = PlatformID.(int)
	delivery.MemberID = MemberID.(int)

	if delivery.ID > 0 {
		err = DB.Updates(&delivery).Error
	} else {
		err = DB.Create(&delivery).Error
	}

	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func MemberDeliveryDelete(c *gin.Context) {
	g := Gin{c}
	var delivery *models.MemberDelivery
	err := c.BindJSON(&delivery)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	MemberID, _ := c.Get("member_id")
	delivery.MemberID = MemberID.(int)

	err = DB.Delete(&delivery).Error

	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	g.Response(http.StatusOK, e.Success, nil)
}
