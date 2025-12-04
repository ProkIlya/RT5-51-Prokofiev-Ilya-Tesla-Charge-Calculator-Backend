package handler

import (
	"net/http"
	"strings"

	"tesla-app/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const jwtPrefix = "Bearer "

// AuthMiddleware проверяет JWT токен
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.Next() // Пропускаем, если нет токена (гость)
			return
		}

		if !strings.HasPrefix(tokenString, jwtPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString = tokenString[len(jwtPrefix):]

		// Проверяем blacklist
		isBlacklisted, err := h.Redis.IsJWTBlacklisted(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}
		if isBlacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
			c.Abort()
			return
		}

		// Парсим токен
		claims := &ds.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.Config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Сохраняем пользователя в контекст
		c.Set("user", &ds.User{
			ID:          claims.UserID,
			Login:       claims.Login,
			IsModerator: claims.IsModerator,
		})

		c.Next()
	}
}

// RequireAuth middleware требует аутентификации
func (h *Handler) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireModerator middleware требует прав модератора
func (h *Handler) RequireModerator() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		currentUser := user.(*ds.User)
		if !currentUser.IsModerator {
			c.JSON(http.StatusForbidden, gin.H{"error": "moderator access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetCurrentUserFromContext получает пользователя из контекста
func (h *Handler) GetCurrentUserFromContext(c *gin.Context) *ds.User {
	user, exists := c.Get("user")
	if !exists || user == nil {
		return nil
	}
	return user.(*ds.User)
}
