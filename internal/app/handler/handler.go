package handler

import (
	"net/http"
	"strconv"
	"strings"

	"tesla-app/internal/app/repository"
	"tesla-app/internal/calculations"
	"tesla-app/internal/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) IndexHandler(c *gin.Context) {
	searchQuery := c.Query("scenario_search")
	scenarios := h.Repository.GetScenarios()
	var filteredScenarios []*models.DrivingScenario

	if searchQuery != "" {
		for _, scenario := range scenarios {
			if strings.Contains(strings.ToLower(scenario.Name), strings.ToLower(searchQuery)) ||
				strings.Contains(strings.ToLower(scenario.Description), strings.ToLower(searchQuery)) {
				filteredScenarios = append(filteredScenarios, scenario)
			}
		}
	} else {
		filteredScenarios = scenarios
	}

	trip := h.Repository.GetTripByID(1)

	c.HTML(http.StatusOK, "drive_menu.html", gin.H{
		"Scenarios":   filteredScenarios,
		"SearchQuery": searchQuery,
		"Trip":        trip,
	})
}

func (h *Handler) ScenarioHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	scenario := h.Repository.GetScenarioByID(id)
	if scenario == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	trip := h.Repository.GetTripByID(1)

	c.HTML(http.StatusOK, "scenario.html", gin.H{
		"Scenario": scenario,
		"Trip":     trip,
	})
}

func (h *Handler) TripHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	trip := h.Repository.GetTripByID(id)
	if trip == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	remainingCharge := calculations.CalculateRemainingCharge(trip)

	c.HTML(http.StatusOK, "trip_calculation.html", gin.H{
		"Trip":            trip,
		"RemainingCharge": remainingCharge,
	})
}
