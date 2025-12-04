package handler

import (
	"net/http"
	"time"

	"tesla-app/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// RegisterUserAPI регистрирует нового пользователя
// @Summary User registration
// @Description Register a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param input body object true "User registration data" example({"login": "string", "password": "string"})
// @Success 201 {object} object "User created"
// @Failure 400 {object} object "Bad request"
// @Router /api/users/register [post]
func (h *Handler) RegisterUserAPI(c *gin.Context) {
	var req struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, не существует ли уже пользователь с таким логином
	existingUser, err := h.Repository.GetUserByLogin(req.Login)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	}

	user := &ds.User{
		Login:       req.Login,
		Password:    req.Password, // Без хеширования по требованию
		IsModerator: false,        // Всегда обычный пользователь
	}

	if err := h.Repository.CreateUser(user); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"login": user.Login,
	})
}

// LoginUserAPI аутентифицирует пользователя
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param input body object true "Login credentials" example({"login": "string", "password": "string"})
// @Success 200 {object} object "Login successful"
// @Failure 401 {object} object "Unauthorized"
// @Router /api/users/login [post]
func (h *Handler) LoginUserAPI(c *gin.Context) {
	var req struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repository.GetUserByLogin(req.Login)
	if err != nil || user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if user == nil || user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Создаем JWT токен
	expirationTime := time.Now().Add(time.Duration(h.Config.JWTExpiresHours) * time.Hour)
	claims := &ds.JWTClaims{
		UserID:      user.ID,
		Login:       user.Login,
		IsModerator: user.IsModerator,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "tesla-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.Config.JWTSecret))
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      tokenString,
		"expires_at": expirationTime,
		"user": gin.H{
			"id":           user.ID,
			"login":        user.Login,
			"is_moderator": user.IsModerator,
		},
	})
}

// LogoutUserAPI выполняет выход пользователя
// @Summary User logout
// @Description Logout user and invalidate token
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object "Logout successful"
// @Router /api/users/logout [post]
func (h *Handler) LogoutUserAPI(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if len(tokenString) > len(jwtPrefix) {
		tokenString = tokenString[len(jwtPrefix):]

		// Добавляем токен в blacklist на оставшееся время
		claims := &ds.JWTClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.Config.JWTSecret), nil
		})

		if err == nil {
			expiresIn := time.Until(time.Unix(claims.ExpiresAt, 0))
			if expiresIn > 0 {
				h.Redis.SetJWTBlacklist(c.Request.Context(), tokenString, expiresIn)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}

// GetUserProfileAPI возвращает профиль текущего пользователя
// @Summary Get user profile
// @Description Get current user profile information
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object "User profile"
// @Router /api/users/profile [get]
func (h *Handler) GetUserProfileAPI(c *gin.Context) {
	user := h.GetCurrentUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"login":        user.Login,
		"is_moderator": user.IsModerator,
	})
}

// UpdateUserProfileAPI обновляет профиль пользователя
// @Summary Update user profile
// @Description Update current user profile information
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body object true "Profile update data" example({"password": "string"})
// @Success 200 {object} object "Profile updated"
// @Router /api/users/profile [put]
func (h *Handler) UpdateUserProfileAPI(c *gin.Context) {
	var req struct {
		Password string `json:"password"`
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

	// Получаем актуальные данные пользователя из базы
	currentUser, err := h.Repository.GetUserByID(user.ID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	if req.Password != "" {
		currentUser.Password = req.Password
		if err := h.Repository.UpdateUser(currentUser); err != nil {
			h.errorHandler(c, http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
