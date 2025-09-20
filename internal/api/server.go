package api

import (
	"tesla-app/internal/app/handler"
	"tesla-app/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func addSecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}

func StartServer() {
	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("Error initializing repository: ", err)
		return
	}

	h := handler.NewHandler(repo)

	router := gin.Default()
	router.Use(addSecurityHeaders())
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	router.GET("/", h.IndexHandler)
	router.GET("/scenario/:id", h.ScenarioHandler)
	router.GET("/trip/trip_:id", h.TripHandler)

	router.Run(":8080")
}
