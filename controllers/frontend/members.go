package frontend

import (
	"eCommerce/internal/auth"
	"eCommerce/pkg/e"
	"encoding/base64"
	"net/http"
	"os"
	"time"

	. "eCommerce/internal/database"
	"eCommerce/models/frontend"
	models "eCommerce/models/frontend"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/liudng/godump"
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

	if req.MemberID == 0 {
		g.Response(http.StatusUnauthorized, e.Unauthorized, err)
		return
	}

	orders, pagination := req.FetchAll()

	for index := range orders {
		orders[index].GetProducts()
	}

	g.PaginationResponse(http.StatusOK, e.Success, orders, pagination)
}

func MemberFetch(c *gin.Context) {
	g := Gin{c}
	if c.Request.Header.Get("Authorization") == "" {
		g.Response(http.StatusOK, e.Success, nil)
		return
	}

	tokenClaims, err := jwt.ParseWithClaims(c.Request.Header.Get("Authorization"), &frontend.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SIGN")), nil
	})

	if err != nil {
		g.Response(http.StatusOK, e.Success, nil)
		return
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*frontend.Claims); ok && tokenClaims.Valid {
			g.Response(http.StatusOK, e.Success, models.Members{Nickname: claims.Nickname})
		} else {
			g.Response(http.StatusOK, e.Success, nil)
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

	issuer := "GiantZongziEC"
	claims := models.Claims{
		MemberID:   member.ID,
		PlatformID: PlatformID.(int),
		Nickname:   member.Nickname,
		StandardClaims: jwt.StandardClaims{
			Issuer: issuer,
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SIGN")))

	/* token := uuid.New().String()

	memberToken := models.MemberToken{
		MemberID:   member.ID,
		Token:      token,
		PlatformID: PlatformID.(int),
	}
	memberToken.CancelOldToken()
	DB.Create(&memberToken) */

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
		godump.Dump(err)
		g.Response(http.StatusOK, e.EmailDuplicate, err)
		return
	}

	token := uuid.New().String()
	memberToken := models.MemberToken{
		MemberID:   member.ID,
		Token:      token,
		PlatformID: PlatformID.(int),
	}
	DB.Create(&memberToken)

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
