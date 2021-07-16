package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	frontend "eCommerce/controllers/frontend"

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
	// r.Use(LoggerToFile())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	frontendApi := router.Group("/api/frontend")
	{
		frontendApi.GET("/test", frontend.Test)
	}

	// var listenTime = 5
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")),
		Handler:        router,
		ReadTimeout:    60,
		WriteTimeout:   60,
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
