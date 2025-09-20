package handler

import (
	"net/http"
	"strconv"

	"tesla-app/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) IndexHandler(c *gin.Context) {
	searchQuery := c.Query("scenario_search")
	var scenarios []*ds.DrivingScenario
	var err error

	if searchQuery != "" {
		scenarios, err = h.Repository.SearchScenarios(searchQuery) // ORM запрос
	} else {
		scenarios, err = h.Repository.GetScenarios() // // ORM запрос
	}

	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Получение заявки пользователя
	trip, err := h.Repository.GetUserDraft(1)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	var cartCount int64
	if trip != nil && trip.Status == "черновик" {
		cartCount = h.Repository.GetTripScenariosCount(trip.ID)
	}

	c.HTML(http.StatusOK, "drive_menu.html", gin.H{
		"Scenarios":   scenarios,
		"SearchQuery": searchQuery,
		"CartCount":   cartCount,
		"HasDraft":    trip != nil && trip.Status == "черновик",
		"CurrentTrip": trip,
	})
}

func (h *Handler) ScenarioHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	scenario, err := h.Repository.GetScenarioByID(uint(id)) // ORM запрос получения одного сценария езды
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	trip, err := h.Repository.GetUserDraft(1)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	var cartCount int64
	if trip != nil && trip.Status == "черновик" {
		cartCount = h.Repository.GetTripScenariosCount(trip.ID)
	}

	c.HTML(http.StatusOK, "scenario.html", gin.H{
		"Scenario":  scenario,
		"CartCount": cartCount,
		"HasDraft":  trip != nil && trip.Status == "черновик",
	})
}

func (h *Handler) AddScenarioToTripHandler(c *gin.Context) {
	userID := uint(1)

	// Получаю или создаю черновик заявки
	trip, err := h.Repository.GetUserDraft(userID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	if trip == nil {
		trip, err = h.Repository.CreateDraft(userID)
		if err != nil {
			h.errorHandler(c, http.StatusInternalServerError, err)
			return
		}
	}

	// Получаю ID сценария
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Получаю значение duration из формы
	durationStr := c.PostForm("duration")
	var duration float64

	if durationStr == "" {
		// Если поле пустое, устанавливаем 0 (без значений по умолчанию)
		duration = 0
	} else {
		duration, err = strconv.ParseFloat(durationStr, 64)
		if err != nil {
			h.errorHandler(c, http.StatusBadRequest, err)
			return
		}
	}

	// Добавляю сценарий в заявку
	err = h.Repository.AddScenarioToTrip(trip.ID, uint(id), duration)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Перенаправляем обратно на главную страницу
	c.Redirect(http.StatusFound, "/")
}
