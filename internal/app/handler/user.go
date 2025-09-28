package handler

import (
	"net/http"

	"tesla-app/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterUserAPI(c *gin.Context) {
	var u ds.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.Repository.CreateUser(&u)
	c.JSON(http.StatusCreated, u)
}

func (h *Handler) GetUserProfileAPI(c *gin.Context) {
	c.JSON(http.StatusOK, GetCurrentUser())
}

func (h *Handler) UpdateUserProfileAPI(c *gin.Context) {
	var req struct {
		Password string `json:"password"`
	}
	c.BindJSON(&req)
	user := GetCurrentUser()
	if req.Password != "" {
		user.Password = req.Password
		h.Repository.UpdateUser(user)
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) LoginUserAPI(c *gin.Context) {
	var cred struct{ Login, Password string }
	c.BindJSON(&cred)
	u, _ := h.Repository.GetUserByLogin(cred.Login)
	if u != nil && u.Password == cred.Password {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid"})
	}
}

func (h *Handler) LogoutUserAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
