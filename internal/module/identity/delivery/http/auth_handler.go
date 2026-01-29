package http

import (
	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/identity/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: username and password required")
		return
	}

	resp, err := h.authUsecase.Login(req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, resp)
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get(auth.UserIDKey)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	user, roles, permissions, err := h.authUsecase.GetCurrentUser(userID.(uint64))
	if err != nil {
		response.InternalError(c, "Failed to get user info")
		return
	}

	response.Success(c, map[string]interface{}{
		"user":        user,
		"roles":       roles,
		"permissions": permissions,
	})
}
