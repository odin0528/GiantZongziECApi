package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"eCommerce/controllers/backend"
	"eCommerce/controllers/frontend"
	auth "eCommerce/internal/auth"
	_ "eCommerce/internal/component"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	env "github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {

	env.Load()

	gin.SetMode("debug")
	gin.ForceConsoleColor()
	router := gin.New()
	router.Use(LoggerToFile())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	frontendApi := router.Group("/api/frontend", auth.GetPlatformID)
	{
		frontendApi.GET("/platform", frontend.PlatformFetch)
		frontendApi.GET("/platform/payment", frontend.PlatformPaymentFetch)

		frontendApi.GET("/pages/:page", frontend.GetPageComponent)

		frontendApi.POST("/products/:layer/:category_id/:page", frontend.GetProductsByCategoryID)
		frontendApi.GET("/categories/:parent_id", frontend.CategoryList)
		frontendApi.GET("/product/:id", frontend.ProductFetch)

		frontendApi.POST("/order/create", frontend.OrderCreate)
		// frontendApi.POST("/order/update", frontend.OrderUpdate)
		// frontendApi.POST("/order/ecpay", frontend.Ecpay)

		frontendApi.GET("/member", frontend.MemberFetch)
		frontendApi.POST("/member/login", frontend.MemberLogin)
		frontendApi.POST("/member/register", frontend.MemberRegister)
		frontendApi.POST("/member/oauth", frontend.MemberOAuth)

		tokenRequired := frontendApi.Use(auth.TokenRequred)
		{
			tokenRequired.POST("/member/orders", frontend.GetMemberOrders)
			tokenRequired.POST("/member/delivery", frontend.GetMemberDelivery)
			tokenRequired.POST("/member/delivery/modify", frontend.MemberDeliveryModify)
			tokenRequired.POST("/member/delivery/delete", frontend.MemberDeliveryDelete)
			// tokenRequired.POST("/member/logout", frontend.MemberLogout)

			tokenRequired.GET("/carts", frontend.CartsFetch)
			tokenRequired.POST("/carts/add", frontend.CartsAddProduct)
			tokenRequired.POST("/carts/update", frontend.CartsUpdate)
			tokenRequired.POST("/carts/remove", frontend.CartsRemoveProduct)
			tokenRequired.POST("/carts/reset", frontend.CartsResetProduct)
		}
	}

	backendApi := router.Group("/api/backend")
	{
		backendApi.POST("/login", backend.Login)
		backendApi.POST("/reset", backend.ResetPassword)
		backendApi.POST("/ecpay/finish", backend.EcpayPaymentFinish)
		backendApi.POST("/ecpay/logistics", backend.EcpayLogisticsNotify)
		backendApi.GET("/ecpay/test", backend.EcpayPaymentTest)

		backendApi.POST("/linepay/finish", backend.LinePayFinish)
		// backendApi.GET("/ecpay/logistics", backend.EcpayLogisticsCreate)

		authRequired := backendApi.Use(auth.AuthRequred)
		{
			authRequired.GET("/pages", backend.GetPagesList)
			authRequired.GET("/pages/:id", backend.GetPageComponent)
			authRequired.POST("/pages/release", backend.PageRelease)
			authRequired.POST("/pages/modify", backend.PageModify)
			// authRequired.POST("/pages/sort", backend.PageSort)
			authRequired.POST("/components/delete", backend.DraftComponentDelete)
			authRequired.POST("/components/create", backend.DraftComponentCreate)
			authRequired.POST("/components/change", backend.DraftComponentChange)
			authRequired.POST("/components/edit", backend.DraftComponentEdit)

			authRequired.GET("/category/list/:parent_id", backend.CategoryList)
			authRequired.GET("/category/fetch/:parent_id", backend.CategoryChildList)
			authRequired.POST("/category/create", backend.CategoryCreate)
			authRequired.POST("/category/modify", backend.CategoryModify)
			authRequired.POST("/category/delete", backend.CategoryDelete)
			authRequired.POST("/category/move", backend.CategoryMove)

			authRequired.POST("/products", backend.ProductList)
			authRequired.GET("/products/:id", backend.ProductFetch)
			authRequired.POST("/products/save", backend.ProductModify)
			authRequired.POST("/products/public", backend.ProductPublic)
			authRequired.POST("/products/delete", backend.ProductDelete)

			authRequired.GET("/orders/untreated", backend.OrderUntreated)
			authRequired.POST("/orders", backend.OrderList)
			authRequired.POST("/order/next", backend.OrderNextStep)
			authRequired.POST("/order/shipment", backend.OrderMakeShipmentNo)
			authRequired.POST("/order/shipment/print", backend.OrderShipmentPrint)

			authRequired.POST("/members/list", backend.MemberList)

			authRequired.POST("/promotions", backend.PromotionList)
			authRequired.POST("/promotions/modify", backend.PromotionModify)

			authRequired.GET("/platform", backend.PlatformFetch)
			authRequired.GET("/platform/payment", backend.PlatformPaymentFetch)
			authRequired.GET("/platform/logistics", backend.PlatformLogisticsFetch)
			authRequired.POST("/platform", backend.PlatformUpdate)
			authRequired.POST("/platform/payment", backend.PlatformPaymentUpdate)
			authRequired.POST("/platform/logistics", backend.PlatformLogisticsUpdate)

			authRequired.GET("/menus", backend.MenuList)
			authRequired.POST("/menus", backend.MenuModify)
			authRequired.POST("/menus/move", backend.MenuMove)
			authRequired.POST("/menus/delete", backend.MenuDelete)

		}
	}

	// var listenTime = 5
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")),
		Handler:        router,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// go runtime.StartListenUSDT()
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Panic: %v", r)
		}
	}()
}

func LoggerToFile() gin.HandlerFunc {

	f, err := os.OpenFile(fmt.Sprintf("logfile_%s%d.log", time.Now().Month(), time.Now().Day()), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		PadLevelText:    true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceQuote:      true,
	})
	log.SetOutput(f)
	log.SetLevel(log.InfoLevel)

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		log.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
