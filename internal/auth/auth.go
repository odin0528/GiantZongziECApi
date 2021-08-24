package auth

import (
	"os"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	// "ec/internal/redis"

	"eCommerce/models/backend"
	"eCommerce/models/frontend"
	"eCommerce/pkg/e"
)

var Salt = []byte{0x47, 0x69, 0x61, 0x6E, 0x74, 0x5a, 0x6F, 0x6E}

func AuthRequred(c *gin.Context) {

	if c.Request.Header.Get("Authorization") == "" {
		Unauthorized(c)
		return
	}

	token := backend.AdminToken{
		Token: c.Request.Header.Get("Authorization"),
	}

	token.Fetch()

	if token.AdminID == 0 {
		Unauthorized(c)
		return
	}

	query := backend.AdminQuery{
		ID: token.AdminID,
	}

	admin := query.Fetch()

	c.Set("platform_id", admin.PlatformID)
}

func TokenRequred(c *gin.Context) {
	if c.Request.Header.Get("Authorization") == "" {
		Unauthorized(c)
		return
	}

	tokenClaims, err := jwt.ParseWithClaims(c.Request.Header.Get("Authorization"), &frontend.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("wJuan"), nil
	})

	if err != nil {
		Unauthorized(c)
		return
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*frontend.Claims); ok && tokenClaims.Valid {
			c.Set("member_id", claims.MemberID)
			c.Set("platform_id", claims.PlatformID)
		} else {
			Unauthorized(c)
			return
		}
	}

	/* PlatformID, _ := c.Get("platform_id")

	token := frontend.MemberToken{
		Token:      c.Request.Header.Get("Authorization"),
		PlatformID: PlatformID.(int),
	}

	token.Fetch()
	if token.MemberID == 0 {
		Unauthorized(c)
		return
	} */

	// c.Set("member_id", token.MemberID)

}

func JwtParser(c *gin.Context) *frontend.Claims {
	tokenClaims, err := jwt.ParseWithClaims(c.Request.Header.Get("Authorization"), &frontend.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SIGN")), nil
	})

	if err != nil {
		return nil
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*frontend.Claims); ok && tokenClaims.Valid {
			return claims
		}
	}

	return nil
}

func GetPlatformID(c *gin.Context) {
	r, _ := regexp.Compile("^([a-zA-Z0-9\\.]*).*$")
	match := r.FindAllStringSubmatch(c.Request.Header["Hostname"][0], 1)

	query := &frontend.PlatformQuery{}
	query.Hostname = match[0][1]
	platform := query.Fetch()

	if platform.ID == 0 {
		Unauthorized(c)
		return
	}

	c.Set("platform_id", platform.ID)
	c.Set("platform", platform)
}

func Unauthorized(c *gin.Context) {
	c.Abort()
	c.JSON(401, gin.H{
		"http_status": 401,
		"code":        401,
		"msg":         e.GetMsg(401),
		"data":        nil,
	})
}
