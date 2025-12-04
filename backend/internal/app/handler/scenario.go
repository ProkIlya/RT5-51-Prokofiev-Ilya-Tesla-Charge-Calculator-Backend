package handler

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"tesla-app/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type ScenarioRequest struct {
	Name             string  `json:"name" binding:"required"`
	Description      string  `json:"description"`
	Type             string  `json:"type" binding:"required,oneof=дорога комфорт"`
	SystemConsuption float64 `json:"system_consumption,omitempty"`
	Speed            float64 `json:"speed,omitempty"`
	AeroCoeff        float64 `json:"aero_coeff,omitempty"`
	RollingCoeff     float64 `json:"rolling_coeff,omitempty"`
}

type ScenarioResponse struct {
	ID               uint    `json:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Status           string  `json:"status"`
	ImageURL         string  `json:"image_url"`
	Type             string  `json:"type"`
	SystemConsuption float64 `json:"system_consumption"`
	Speed            float64 `json:"speed"`
	AeroCoeff        float64 `json:"aero_coeff"`
	RollingCoeff     float64 `json:"rolling_coeff"`
}

// GetScenariosAPI возвращает список сценариев
// @Summary Get scenarios list
// @Description Get all driving scenarios with optional filtering
// @Tags Scenarios
// @Accept json
// @Produce json
// @Param name query string false "Filter by name"
// @Param type query string false "Filter by type"
// @Success 200 {array} ScenarioResponse
// @Failure 500 {object} object "Internal server error"
// @Router /api/scenarios [get]
func (h *Handler) GetScenariosAPI(c *gin.Context) {
	nameFilter := c.Query("name")
	typeFilter := c.Query("type")
	scenarios, err := h.Repository.GetScenariosWithFilters(nameFilter, typeFilter)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}
	resp := make([]ScenarioResponse, len(scenarios))
	for i, s := range scenarios {
		resp[i] = h.scenarioToResponse(s)
	}
	c.JSON(http.StatusOK, resp)
}

// GetScenarioAPI возвращает сценарий по ID
// @Summary Get scenario by ID
// @Description Get driving scenario details by ID
// @Tags Scenarios
// @Accept json
// @Produce json
// @Param id path int true "Scenario ID"
// @Success 200 {object} ScenarioResponse
// @Failure 400 {object} object "Invalid ID"
// @Failure 404 {object} object "Scenario not found"
// @Router /api/scenarios/{id} [get]
func (h *Handler) GetScenarioAPI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	s, err := h.Repository.GetScenarioByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, h.scenarioToResponse(s))
}

// CreateScenarioAPI создает новый сценарий
// @Summary Create scenario
// @Description Create new driving scenario (moderator only)
// @Tags Scenarios
// @Accept json
// @Produce json
// @Param input body ScenarioRequest true "Scenario data"
// @Success 201 {object} ScenarioResponse
// @Failure 400 {object} object "Bad request"
// @Failure 500 {object} object "Internal server error"
// @Security BearerAuth
// @Router /api/scenarios [post]
func (h *Handler) CreateScenarioAPI(c *gin.Context) {
	var req ScenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	scenario := &ds.DrivingScenario{
		Name:             req.Name,
		Description:      req.Description,
		Status:           "действует",
		Type:             req.Type,
		SystemConsuption: req.SystemConsuption,
		Speed:            req.Speed,
		AeroCoeff:        req.AeroCoeff,
		RollingCoeff:     req.RollingCoeff,
	}
	if err := h.Repository.CreateScenario(scenario); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, h.scenarioToResponse(scenario))
}

// UpdateScenarioAPI обновляет сценарий
// @Summary Update scenario
// @Description Update driving scenario (moderator only)
// @Tags Scenarios
// @Accept json
// @Produce json
// @Param id path int true "Scenario ID"
// @Param input body ScenarioRequest true "Scenario data"
// @Success 200 {object} ScenarioResponse
// @Failure 400 {object} object "Bad request"
// @Failure 404 {object} object "Scenario not found"
// @Failure 500 {object} object "Internal server error"
// @Security BearerAuth
// @Router /api/scenarios/{id} [put]
func (h *Handler) UpdateScenarioAPI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req ScenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	scenario, err := h.Repository.GetScenarioByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	scenario.Name = req.Name
	scenario.Description = req.Description
	scenario.Type = req.Type
	scenario.SystemConsuption = req.SystemConsuption
	scenario.Speed = req.Speed
	scenario.AeroCoeff = req.AeroCoeff
	scenario.RollingCoeff = req.RollingCoeff
	if err := h.Repository.UpdateScenario(scenario); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, h.scenarioToResponse(scenario))
}

// DeleteScenarioAPI удаляет сценарий
// @Summary Delete scenario
// @Description Delete driving scenario (moderator only)
// @Tags Scenarios
// @Accept json
// @Produce json
// @Param id path int true "Scenario ID"
// @Success 200 {object} object "Deleted successfully"
// @Failure 400 {object} object "Invalid ID"
// @Failure 404 {object} object "Scenario not found"
// @Failure 500 {object} object "Internal server error"
// @Security BearerAuth
// @Router /api/scenarios/{id} [delete]
func (h *Handler) DeleteScenarioAPI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	scenario, err := h.Repository.GetScenarioByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if scenario.ImageURL != "" {
		h.deleteImageFromMinio(scenario.ImageURL)
	}
	if err := h.Repository.DeleteScenario(uint(id)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// AddScenarioToTripAPI добавляет сценарий в заявку
// @Summary Add scenario to trip
// @Description Add driving scenario to user's draft trip
// @Tags Scenarios
// @Accept json
// @Produce json
// @Param id path int true "Scenario ID"
// @Param input body object true "Duration data" example({"duration": 2.5})
// @Success 200 {object} object "Scenario added to trip"
// @Failure 401 {object} object "Unauthorized"
// @Failure 500 {object} object "Internal server error"
// @Security BearerAuth
// @Router /api/scenarios/{id}/add-to-trip [post]
func (h *Handler) AddScenarioToTripAPI(c *gin.Context) {
	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	scenarioID, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Duration float64 `json:"duration"`
	}
	c.ShouldBindJSON(&req)
	trip, _ := h.Repository.GetUserDraft(user.ID)
	if trip == nil {
		trip, _ = h.Repository.CreateDraft(user.ID)
	}
	h.Repository.AddScenarioToTrip(trip.ID, uint(scenarioID), req.Duration)
	c.JSON(http.StatusOK, gin.H{"message": "added"})
}

// UploadScenarioImageAPI загружает изображение для сценария
// @Summary Upload scenario image
// @Description Upload image for driving scenario (moderator only)
// @Tags Scenarios
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Scenario ID"
// @Param image formData file true "Image file"
// @Success 200 {object} object "Image uploaded"
// @Failure 400 {object} object "Bad request"
// @Failure 500 {object} object "Internal server error"
// @Security BearerAuth
// @Router /api/scenarios/{id}/scenarioimage [post]
func (h *Handler) UploadScenarioImageAPI(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	scenario, _ := h.Repository.GetScenarioByID(uint(id))
	file, hdr, _ := c.Request.FormFile("image")
	defer file.Close()
	if scenario.ImageURL != "" {
		h.deleteImageFromMinio(scenario.ImageURL)
	}
	name := fmt.Sprintf("%d_%s%s", time.Now().Unix(), strconv.Itoa(id), filepath.Ext(hdr.Filename))
	url, _ := h.uploadImageToMinio(file, name, hdr.Size)
	scenario.ImageURL = url
	h.Repository.UpdateScenario(scenario)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *Handler) scenarioToResponse(scenario *ds.DrivingScenario) ScenarioResponse {
	return ScenarioResponse{
		ID:               scenario.ID,
		Name:             scenario.Name,
		Description:      scenario.Description,
		Status:           scenario.Status,
		ImageURL:         scenario.ImageURL,
		Type:             scenario.Type,
		SystemConsuption: scenario.SystemConsuption,
		Speed:            scenario.Speed,
		AeroCoeff:        scenario.AeroCoeff,
		RollingCoeff:     scenario.RollingCoeff,
	}
}

func (h *Handler) uploadImageToMinio(file multipart.File, fileName string, fileSize int64) (string, error) {
	bucket := "tesla-images"
	ctx := context.Background()

	// Используем h.MinioClient вместо глобальной переменной
	if exists, _ := h.MinioClient.BucketExists(ctx, bucket); !exists {
		if err := h.MinioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return "", err
		}
	}

	if _, err := h.MinioClient.PutObject(ctx, bucket, fileName, file, fileSize, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	}); err != nil {
		return "", err
	}

	return fmt.Sprintf("http://localhost:9000/%s/%s", bucket, fileName), nil
}

func (h *Handler) deleteImageFromMinio(imageURL string) {
	parts := strings.Split(imageURL, "/")
	if len(parts) < 2 || h.MinioClient == nil {
		return
	}
	bucket := parts[len(parts)-2]
	object := parts[len(parts)-1]
	ctx := context.Background()
	h.MinioClient.RemoveObject(ctx, bucket, object, minio.RemoveObjectOptions{})
}

func (h *Handler) generateImageFileName(original string) string {
	ext := filepath.Ext(original)
	name := strings.TrimSuffix(original, ext)
	name = strings.ReplaceAll(name, " ", "_")
	return fmt.Sprintf("%s_%d%s", name, time.Now().Unix(), ext)
}
