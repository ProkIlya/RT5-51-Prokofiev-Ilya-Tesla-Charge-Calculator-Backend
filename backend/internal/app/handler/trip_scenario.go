package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RemoveScenarioFromTripAPI удаляет сценарий из заявки
// @Summary Remove scenario from trip
// @Description Remove scenario from trip (only for draft trips by owner)
// @Tags TripScenarios
// @Accept json
// @Produce json
// @Param trip_id path int true "Trip ID"
// @Param scenario_id path int true "Scenario ID"
// @Success 200 {object} object "Scenario removed from trip"
// @Failure 400 {object} object "Bad request"
// @Failure 403 {object} object "Access denied"
// @Failure 404 {object} object "Trip not found"
// @Security BearerAuth
// @Router /api/trips/{trip_id}/scenarios/{scenario_id} [delete]
func (h *Handler) RemoveScenarioFromTripAPI(c *gin.Context) {
	tripID, err := strconv.Atoi(c.Param("trip_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip id"})
		return
	}

	scenarioID, err := strconv.Atoi(c.Param("scenario_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scenario id"})
		return
	}

	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	trip, err := h.Repository.GetTripByID(uint(tripID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trip not found"})
		return
	}

	// Проверяем что пользователь является создателем и заявка в черновике
	if trip.CreatorID != user.ID || trip.Status != "черновик" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := h.Repository.RemoveScenarioFromTrip(uint(tripID), uint(scenarioID)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "scenario removed from trip"})
}

// UpdateTripScenarioAPI обновляет сценарий в заявке
// @Summary Update trip scenario
// @Description Update scenario duration in trip (only for draft trips by owner)
// @Tags TripScenarios
// @Accept json
// @Produce json
// @Param trip_id path int true "Trip ID"
// @Param scenario_id path int true "Scenario ID"
// @Param input body object true "Duration data" example({"duration": 3.0})
// @Success 200 {object} object "Trip scenario updated"
// @Failure 400 {object} object "Bad request"
// @Failure 403 {object} object "Access denied"
// @Failure 404 {object} object "Trip not found"
// @Security BearerAuth
// @Router /api/trips/{trip_id}/scenarios/{scenario_id} [put]
func (h *Handler) UpdateTripScenarioAPI(c *gin.Context) {
	tripID, err := strconv.Atoi(c.Param("trip_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip id"})
		return
	}

	scenarioID, err := strconv.Atoi(c.Param("scenario_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scenario id"})
		return
	}

	var req struct {
		Duration float64 `json:"duration" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	trip, err := h.Repository.GetTripByID(uint(tripID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trip not found"})
		return
	}

	if trip.CreatorID != user.ID || trip.Status != "черновик" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := h.Repository.UpdateTripScenario(uint(tripID), uint(scenarioID), req.Duration); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trip scenario updated"})
}

// POST /api/trips/:trip_id/scenarios/:scenario_id/delete - удаление сценария из заявки
func (h *Handler) RemoveScenarioFromTripPostAPI(c *gin.Context) {
	tripID, err := strconv.Atoi(c.Param("trip_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip id"})
		return
	}

	scenarioID, err := strconv.Atoi(c.Param("scenario_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scenario id"})
		return
	}

	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	trip, err := h.Repository.GetTripByID(uint(tripID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trip not found"})
		return
	}

	// Проверяем что пользователь является создателем и заявка в черновике
	if trip.CreatorID != user.ID || trip.Status != "черновик" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := h.Repository.RemoveScenarioFromTrip(uint(tripID), uint(scenarioID)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Перенаправляем обратно на страницу заявки
	c.Redirect(http.StatusSeeOther, "/trip/trip_"+strconv.Itoa(tripID))
}
