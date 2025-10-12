package handler

import (
	"net/http"
	"strconv"
	"tesla-app/internal/app/ds"
	"tesla-app/internal/calculations"
	"time"

	"github.com/gin-gonic/gin"
)

type RussianTime time.Time

// MarshalJSON реализует кастомную сериализацию JSON
func (rt RussianTime) MarshalJSON() ([]byte, error) {
	t := time.Time(rt)
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + t.Format("02.01.2006") + `"`), nil
}

// UnmarshalJSON реализует кастомную десериализацию JSON
func (rt *RussianTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*rt = RussianTime(time.Time{})
		return nil
	}
	t, err := time.Parse(`"`+"02.01.2006"+`"`, string(data))
	if err != nil {
		return err
	}
	*rt = RussianTime(t)
	return nil
}

// DTO для корзины
type CartResponse struct {
	TripID uint  `json:"trip_id"`
	Count  int64 `json:"count"`
}

// DTO для ответа по заявке
type TripResponse struct {
	ID              uint                   `json:"id"`
	Status          string                 `json:"status"`
	CreatorLogin    string                 `json:"creator_login"`
	ModeratorLogin  *string                `json:"moderator_login,omitempty"`
	StartCharge     *float64               `json:"start_charge"`
	RemainingCharge *float64               `json:"remaining_charge"` // убрал omitempty
	CreatedAt       RussianTime            `json:"created_at"`       // Изменено на RussianTime
	SubmittedAt     *RussianTime           `json:"submitted_at,omitempty"`
	CompletedAt     *RussianTime           `json:"completed_at,omitempty"`
	Scenarios       []TripScenarioResponse `json:"scenarios,omitempty"`
}

type TripScenarioResponse struct {
	ScenarioID uint             `json:"scenario_id"`
	Duration   *float64         `json:"duration"`
	Scenario   ScenarioResponse `json:"scenario"`
}

// Преобразование trip и scenarios в TripResponse
func (h *Handler) tripToResponse(trip *ds.TripApplication, ts []TripScenarioResponse) TripResponse {
	resp := TripResponse{
		ID:              trip.ID,
		Status:          trip.Status,
		CreatorLogin:    trip.Creator.Login,
		StartCharge:     trip.StartCharge,
		RemainingCharge: trip.RemainingCharge,
		CreatedAt:       RussianTime(trip.CreatedAt), // Конвертация
		Scenarios:       ts,
	}

	// Конвертация для указателей
	if trip.SubmittedAt != nil {
		submitted := RussianTime(*trip.SubmittedAt)
		resp.SubmittedAt = &submitted
	}
	if trip.CompletedAt != nil {
		completed := RussianTime(*trip.CompletedAt)
		resp.CompletedAt = &completed
	}

	if trip.ModeratorID != nil {
		resp.ModeratorLogin = &trip.Moderator.Login
	}
	return resp
}

// Вспомогательная функция конвертации []ds.TripScenario в []TripScenarioResponse
func (h *Handler) convert(ts []ds.TripScenario) []TripScenarioResponse {
	out := make([]TripScenarioResponse, len(ts))
	for i, t := range ts {
		out[i] = TripScenarioResponse{
			ScenarioID: t.DrivingScenarioID,
			Duration:   t.Duration,
			Scenario:   h.scenarioToResponse(&t.DrivingScenario),
		}
	}
	return out
}

func (h *Handler) GetScenariosCartAPI(c *gin.Context) {
	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		// Для гостя возвращаем 0
		c.JSON(http.StatusOK, CartResponse{TripID: 0, Count: 0})
		return
	}
	draft, _ := h.Repository.GetUserDraft(user.ID)
	if draft == nil {
		c.JSON(http.StatusOK, CartResponse{TripID: 0, Count: 0})
		return
	}
	cnt := h.Repository.GetTripScenariosCount(draft.ID)
	c.JSON(http.StatusOK, CartResponse{TripID: draft.ID, Count: cnt})
}

// GetTripsAPI возвращает список заявок
// @Summary Get trips list
// @Description Get trips with filtering. For moderators - all trips, for users - only their trips
// @Tags Trips
// @Accept json
// @Produce json
// @Param status query string false "Status filter"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Security BearerAuth
// @Success 200 {array} TripResponse
// @Failure 401 {object} object "Unauthorized"
// @Failure 500 {object} object "Internal server error"
// @Router /api/trips [get]
func (h *Handler) GetTripsAPI(c *gin.Context) {
	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	status := c.Query("status")
	var from, to *time.Time
	if df := c.Query("date_from"); df != "" {
		d, _ := time.Parse("2006-01-02", df)
		from = &d
	}
	if dt := c.Query("date_to"); dt != "" {
		d, _ := time.Parse("2006-01-02", dt)
		to = &d
	}

	var trips []*ds.TripApplication
	var err error

	if user.IsModerator {
		// Модератор - все заявки (userID = nil)
		trips, err = h.Repository.GetTripsWithFilters(status, from, to, nil)
	} else {
		// Обычный пользователь - только свои заявки
		trips, err = h.Repository.GetTripsWithFilters(status, from, to, &user.ID)
	}

	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	resp := make([]TripResponse, len(trips))
	for i, t := range trips {
		resp[i] = h.tripToResponse(t, nil)
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetTripAPI(c *gin.Context) {
	tripID, _ := strconv.Atoi(c.Param("trip_id"))
	trip, err := h.Repository.GetTripByIDWithDetails(uint(tripID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trip not found"})
		return
	}

	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Проверяем права доступа: модератор или создатель заявки
	if !user.IsModerator && trip.CreatorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	ts, _ := h.Repository.GetTripScenarios(uint(tripID))
	resp := h.tripToResponse(trip, h.convert(ts))
	c.JSON(http.StatusOK, resp)
}

// 11. PUT /api/trips/:id
func (h *Handler) UpdateTripAPI(c *gin.Context) {
	tripID, err := strconv.Atoi(c.Param("trip_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req struct {
		StartCharge float64 `json:"start_charge"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	trip, err := h.Repository.GetTripByID(uint(tripID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	if trip.CreatorID != user.ID || trip.Status != "черновик" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	trip.StartCharge = &req.StartCharge
	h.Repository.UpdateTrip(trip)
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// 14. DELETE /api/trips/:id
func (h *Handler) DeleteTripAPI(c *gin.Context) {
	tripID, err := strconv.Atoi(c.Param("trip_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	trip, err := h.Repository.GetTripByID(uint(tripID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Проверяем права: только создатель может удалить свою заявку
	if trip.CreatorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := h.Repository.DeleteTrip(uint(tripID)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// PUT /api/trips/:id/reject - отклонение заявки модератором
func (h *Handler) RejectTripAPI(c *gin.Context) {
	tripID, err := strconv.Atoi(c.Param("trip_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	trip, err := h.Repository.GetTripByID(uint(tripID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	if !user.IsModerator || trip.Status != "сформирован" {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot reject trip"})
		return
	}

	now := time.Now()
	trip.Status = "отклонён"
	trip.ModeratorID = &user.ID
	trip.CompletedAt = &now

	if err := h.Repository.UpdateTrip(trip); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trip rejected"})
}

func (h *Handler) SubmitTripAPI(c *gin.Context) {
	tripID, err := strconv.Atoi(c.Param("trip_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	trip, err := h.Repository.GetTripByID(uint(tripID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	if trip.CreatorID != user.ID || trip.Status != "черновик" {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot submit trip"})
		return
	}

	// Проверяем что есть сценарии в заявке
	if h.Repository.GetTripScenariosCount(trip.ID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trip must have scenarios"})
		return
	}

	now := time.Now()
	trip.Status = "сформирован"
	trip.SubmittedAt = &now

	if err := h.Repository.UpdateTrip(trip); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trip submitted"})
}

func (h *Handler) CompleteTripAPI(c *gin.Context) {
	tripID, err := strconv.Atoi(c.Param("trip_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	trip, err := h.Repository.GetTripByID(uint(tripID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	if !user.IsModerator || trip.Status != "сформирован" {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot complete trip"})
		return
	}

	now := time.Now()
	trip.Status = "завершён"
	trip.ModeratorID = &user.ID
	trip.CompletedAt = &now

	// Расчет оставшегося заряда
	scenarios, err := h.Repository.GetTripScenarios(trip.ID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	remainingCharge := calculations.CalculateRemainingCharge(trip, scenarios)
	trip.RemainingCharge = &remainingCharge

	if err := h.Repository.UpdateTrip(trip); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "trip completed",
		"remaining_charge": remainingCharge,
	})
}

// ReviewTripAPI завершает или отклоняет заявку
// @Summary Review trip
// @Description Complete or reject trip (moderator only)
// @Tags Trips
// @Accept json
// @Produce json
// @Param trip_id path int true "Trip ID"
// @Param input body object true "Review action" example({"action":"complete"})
// @Security BearerAuth
// @Success 200 {object} object "Review successful"
// @Failure 400 {object} object "Bad request"
// @Failure 403 {object} object "Forbidden"
// @Failure 404 {object} object "Not found"
// @Router /api/trips/{trip_id}/reviewtrip [put]
func (h *Handler) ReviewTripAPI(c *gin.Context) {
	tripID, err := strconv.Atoi(c.Param("trip_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Action string `json:"action" binding:"required,oneof=complete reject"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trip, err := h.Repository.GetTripByID(uint(tripID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	if !user.IsModerator || trip.Status != "сформирован" {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot review trip"})
		return
	}

	now := time.Now()
	trip.ModeratorID = &user.ID
	trip.CompletedAt = &now

	if req.Action == "complete" {
		trip.Status = "завершён"
		// Расчет оставшегося заряда
		scenarios, err := h.Repository.GetTripScenarios(trip.ID)
		if err != nil {
			h.errorHandler(c, http.StatusInternalServerError, err)
			return
		}
		remainingCharge := calculations.CalculateRemainingCharge(trip, scenarios)
		trip.RemainingCharge = &remainingCharge
	} else if req.Action == "reject" {
		trip.Status = "отклонён"
	}

	if err := h.Repository.UpdateTrip(trip); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "trip " + req.Action + "ed",
		"remaining_charge": trip.RemainingCharge,
	})
}

// POST /api/trips - создание новой заявки
func (h *Handler) CreateTripAPI(c *gin.Context) {
	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req struct {
		StartCharge float64 `json:"start_charge" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trip := &ds.TripApplication{
		CreatorID:   user.ID,
		Status:      "черновик",
		StartCharge: &req.StartCharge,
	}

	if err := h.Repository.CreateTrip(trip); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           trip.ID,
		"status":       trip.Status,
		"start_charge": trip.StartCharge,
		"created_at":   trip.CreatedAt,
	})
}
