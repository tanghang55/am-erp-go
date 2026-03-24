package http

import (
	"errors"
	"strconv"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/identity/domain"
	"am-erp-go/internal/module/identity/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	params := &domain.UserListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
		Status:   c.Query("status"),
		Keyword:  c.Query("keyword"),
	}

	list, total, err := h.userUsecase.ListUsers(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, list, total, params.Page, params.PageSize)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	user, roles, permissions, err := h.userUsecase.GetUserDetail(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"user":        user,
		"roles":       roles,
		"permissions": permissions,
	})
}

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RealName string `json:"real_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.userUsecase.CreateUser(&usecase.CreateUserInput{
		Username: req.Username,
		Password: req.Password,
		RealName: req.RealName,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, user)
}

type updateUserRequest struct {
	Password *string `json:"password"`
	RealName *string `json:"real_name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Status   *string `json:"status"`
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.userUsecase.UpdateUser(id, &usecase.UpdateUserInput{
		Password: req.Password,
		RealName: req.RealName,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if currentUserID, ok := getUserIDFromContext(c); ok && currentUserID == id {
		response.BadRequest(c, "cannot disable current user")
		return
	}

	if err := h.userUsecase.DisableUser(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

type assignRolesRequest struct {
	RoleIDs []uint64 `json:"role_ids" binding:"required"`
}

func (h *UserHandler) AssignUserRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req assignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userUsecase.AssignRoles(id, req.RoleIDs); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) ListRoles(c *gin.Context) {
	roles, err := h.userUsecase.ListRoles()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, roles)
}

func (h *UserHandler) ListPermissions(c *gin.Context) {
	permissions, err := h.userUsecase.ListPermissions()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, permissions)
}

func parseIntOrDefault(s string, defaultValue int) int {
	value, err := strconv.Atoi(s)
	if err != nil || value <= 0 {
		return defaultValue
	}
	return value
}

func getUserIDFromContext(c *gin.Context) (uint64, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint64)
	if !ok {
		return 0, false
	}
	return id, true
}
