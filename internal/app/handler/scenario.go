package handler

import (
	"net/http"
	"strconv"

	"tesla-app/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) IndexHandler(c *gin.Context) {
	searchQuery := c.Query("search")
	var scenarios []*ds.DrivingScenario
	var err error

	if searchQuery != "" {
		scenarios, err = h.Repository.SearchScenarios(searchQuery)
	} else {
		scenarios, err = h.Repository.GetScenarios()
	}

	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Получаем заявку пользователя (пока используем ID 1 для демонстрации)
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

	scenario, err := h.Repository.GetScenarioByID(uint(id))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Получаем заявку пользователя (пока используем ID 1 для демонстрации)
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
	// Получаем ID пользователя (пока используем ID 1 для демонстрации)
	userID := uint(1)

	// Получаем или создаем черновик заявки
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

	// Получаем ID сценария
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Получаем сценарий для определения типа
	scenario, err := h.Repository.GetScenarioByID(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Устанавливаем значение по умолчанию в зависимости от типа сценария
	var defaultValue float64
	if scenario.Type == "дорога" {
		defaultValue = 50
	} else {
		defaultValue = 1
	}

	// Получаем значение из формы или используем значение по умолчанию
	valueStr := c.PostForm("value")
	if valueStr == "" {
		valueStr = strconv.FormatFloat(defaultValue, 'f', -1, 64)
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	// Добавляем сценарий в заявку
	err = h.Repository.AddScenarioToTrip(trip.ID, uint(id), value)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Перенаправляем обратно на главную страницу
	c.Redirect(http.StatusFound, "/")
}
