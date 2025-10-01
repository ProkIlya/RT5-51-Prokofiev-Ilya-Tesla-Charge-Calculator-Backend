package handler

import (
	"tesla-app/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type Handler struct {
	Repository  *repository.Repository
	MinioClient *minio.Client
}

func NewHandler(r *repository.Repository, mc *minio.Client) *Handler {
	return &Handler{
		Repository:  r,
		MinioClient: mc,
	}
}

func (h *Handler) errorHandler(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api")

	scenarios := api.Group("/scenarios")
	{
		scenarios.GET("", h.GetScenariosAPI)
		scenarios.GET("/:id", h.GetScenarioAPI)
		scenarios.POST("", h.CreateScenarioAPI)
		scenarios.PUT("/:id", h.UpdateScenarioAPI)
		scenarios.DELETE("/:id", h.DeleteScenarioAPI)
		scenarios.POST("/:id/add-to-trip", h.AddScenarioToTripAPI)
		scenarios.POST("/:id/scenarioimage", h.UploadScenarioImageAPI)
	}

	trips := api.Group("/trips")
	{
		trips.GET("/scenarioscart", h.GetScenariosCartAPI)
		trips.GET("", h.GetTripsAPI)
		//trips.POST("", h.CreateTripAPI)
		trips.GET("/:trip_id", h.GetTripAPI)
		trips.PUT("/:trip_id", h.UpdateTripAPI)
		trips.PUT("/:trip_id/submittrip", h.SubmitTripAPI)
		trips.PUT("/:trip_id/reviewtrip", h.ReviewTripAPI)
		trips.DELETE("/:trip_id", h.DeleteTripAPI)

		tripScenarios := trips.Group("/:trip_id/scenarios")
		{
			tripScenarios.DELETE("/:scenario_id", h.RemoveScenarioFromTripAPI)
			//tripScenarios.POST("/:scenario_id/delete", h.RemoveScenarioFromTripPostAPI)
			tripScenarios.PUT("/:scenario_id", h.UpdateTripScenarioAPI)
		}
	}

	users := api.Group("/users")
	{
		users.POST("/register", h.RegisterUserAPI)
		users.GET("/profile", h.GetUserProfileAPI)
		users.PUT("/profile", h.UpdateUserProfileAPI)
		users.POST("/login", h.LoginUserAPI)
		users.POST("/logout", h.LogoutUserAPI)
	}
}
