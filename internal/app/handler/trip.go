package handler

import (
	"net/http"
	"strconv"

	"tesla-app/internal/calculations"

	"github.com/gin-gonic/gin"
)

func (h *Handler) TripHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	trip, err := h.Repository.GetTripByID(uint(id)) // ORM запрос
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Проверяю, что заявка не удалена
	if trip.Status == "удалён" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Получаю сценарии для этой заявки
	tripScenarios, err := h.Repository.GetTripScenarios(uint(id)) // ORM запрос получения заявки по id
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Получаю количество услуг в корзине для текущего пользователя
	userTrip, err := h.Repository.GetUserDraft(1) // ID 1 для демонстрации
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	var cartCount int64
	var hasDraft bool
	if userTrip != nil && userTrip.Status == "черновик" {
		cartCount = h.Repository.GetTripScenariosCount(userTrip.ID)
		hasDraft = true
	}

	// Расчет остатка заряда
	remainingCharge := calculations.CalculateRemainingCharge(trip, tripScenarios)

	c.HTML(http.StatusOK, "trip_calculation.html", gin.H{
		"Trip":            trip,
		"TripScenarios":   tripScenarios,
		"RemainingCharge": remainingCharge,
		"CartCount":       cartCount,
		"HasDraft":        hasDraft,
	})
}

func (h *Handler) DeleteTripHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Удаляем заявку через прямой SQL (без ORM)
	err = h.Repository.DeleteTrip(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Перенаправляем на главную страницу
	c.Redirect(http.StatusFound, "/")
}
