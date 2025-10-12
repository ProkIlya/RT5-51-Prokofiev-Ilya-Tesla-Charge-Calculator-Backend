package handler

import (
	"tesla-app/internal/app/config"
	"tesla-app/internal/app/redis"
	"tesla-app/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type Handler struct {
	Repository  *repository.Repository
	MinioClient *minio.Client
	Config      *config.Config
	Redis       *redis.Client
}

func NewHandler(r *repository.Repository, mc *minio.Client, cfg *config.Config, redisClient *redis.Client) *Handler {
	return &Handler{
		Repository:  r,
		MinioClient: mc,
		Config:      cfg,
		Redis:       redisClient,
	}
}

func (h *Handler) errorHandler(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	// Глобальный middleware для аутентификации
	router.Use(h.AuthMiddleware())

	api := router.Group("/api")

	// Публичные маршруты (доступны всем)
	public := api.Group("")
	{
		// Аутентификация (доступна только гостям)
		public.POST("/users/register", h.RegisterUserAPI)
		public.POST("/users/login", h.LoginUserAPI)
		// Просмотр сценариев (доступен всем)
		public.GET("/scenarios", h.GetScenariosAPI)
		public.GET("/scenarios/:id", h.GetScenarioAPI)
		// Корзина (доступна всем, но для гостя возвращает 0,0)
		public.GET("/trips/scenarioscart", h.GetScenariosCartAPI)
	}

	// Защищенные маршруты (требуют аутентификации)
	protected := api.Group("")
	protected.Use(h.RequireAuth())
	{
		// Пользовательские маршруты
		users := protected.Group("/users")
		{
			users.GET("/profile", h.GetUserProfileAPI)
			users.PUT("/profile", h.UpdateUserProfileAPI)
			users.POST("/logout", h.LogoutUserAPI)
		}

		// Маршруты заявок (для всех аутентифицированных)
		trips := protected.Group("/trips")
		{
			//trips.GET("/scenarioscart", h.GetScenariosCartAPI)
			trips.GET("", h.GetTripsAPI)
			trips.GET("/:trip_id", h.GetTripAPI)
			trips.PUT("/:trip_id", h.UpdateTripAPI)
			trips.PUT("/:trip_id/submittrip", h.SubmitTripAPI)
			trips.DELETE("/:trip_id", h.DeleteTripAPI)

			tripScenarios := trips.Group("/:trip_id/scenarios")
			{
				tripScenarios.DELETE("/:scenario_id", h.RemoveScenarioFromTripAPI)
				tripScenarios.PUT("/:scenario_id", h.UpdateTripScenarioAPI)
			}
		}

		// Маршруты сценариев (чтение для всех, изменение для модераторов)
		scenarios := protected.Group("/scenarios")
		{
			scenarios.POST("", h.RequireModerator(), h.CreateScenarioAPI)
			scenarios.PUT("/:id", h.RequireModerator(), h.UpdateScenarioAPI)
			scenarios.DELETE("/:id", h.RequireModerator(), h.DeleteScenarioAPI)
			scenarios.POST("/:id/add-to-trip", h.AddScenarioToTripAPI)
			scenarios.POST("/:id/scenarioimage", h.RequireModerator(), h.UploadScenarioImageAPI)
		}

		// Маршруты только для модераторов
		moderator := protected.Group("")
		moderator.Use(h.RequireModerator())
		{
			moderator.PUT("/trips/:trip_id/reviewtrip", h.ReviewTripAPI)
		}
	}
}
