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
	// router.Use(LoggerToFile())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	frontendApi := router.Group("/api/frontend", auth.GetCustomerID)
	{
		frontendApi.GET("/pages/:page", frontend.GetPageComponent)
		// frontendApi.GET("/products/:layer/:category_id/:page", frontend.GetProductsByCategoryID)
		frontendApi.POST("/products/:layer/:category_id/:page", frontend.GetProductsByCategoryID)
		frontendApi.GET("/categories/:parent_id", frontend.CategoryList)
	}

	backendApi := router.Group("/api/backend", auth.AuthRequred)
	{
		backendApi.GET("/pages", backend.GetPagesList)
		backendApi.GET("/pages/:page_id", backend.GetPageComponent)
		backendApi.POST("/components/delete", backend.DraftComponentDelete)
		backendApi.POST("/components/create", backend.DraftComponentCreate)
		backendApi.POST("/components/change", backend.DraftComponentChange)
		backendApi.POST("/components/edit", backend.DraftComponentEdit)

		backendApi.GET("/category/list/:parent_id", backend.CategoryList)
		backendApi.GET("/category/fetch/:parent_id", backend.CategoryChildList)
		backendApi.POST("/category/create", backend.CategoryCreate)
		backendApi.POST("/category/modify", backend.CategoryModify)
		backendApi.POST("/category/delete", backend.CategoryDelete)
		backendApi.POST("/category/move", backend.CategoryMove)

		backendApi.POST("/products", backend.ProductList)
		backendApi.GET("/products/:id", backend.ProductFetch)
		backendApi.POST("/products/save", backend.ProductModify)
		backendApi.POST("/products/public", backend.ProductPublic)
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

	f, err := os.OpenFile("testlogfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
