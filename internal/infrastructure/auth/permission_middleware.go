package auth

import (
	"strings"

	"am-erp-go/internal/infrastructure/response"
	identityDomain "am-erp-go/internal/module/identity/domain"

	"github.com/gin-gonic/gin"
)

type PermissionRepository interface {
	GetUserRoles(userID uint64) ([]identityDomain.Role, error)
	GetUserPermissions(userID uint64) ([]identityDomain.Permission, error)
}

func RequirePermission(repo PermissionRepository, codes ...string) gin.HandlerFunc {
	normalized := make([]string, 0, len(codes))
	for _, code := range codes {
		code = strings.TrimSpace(code)
		if code != "" {
			normalized = append(normalized, code)
		}
	}

	return func(c *gin.Context) {
		if len(normalized) == 0 || repo == nil {
			c.Next()
			return
		}

		userIDValue, exists := c.Get(UserIDKey)
		if !exists {
			response.Unauthorized(c, "user identity missing")
			c.Abort()
			return
		}

		userID, ok := userIDValue.(uint64)
		if !ok {
			response.Unauthorized(c, "invalid user identity")
			c.Abort()
			return
		}

		roles, err := repo.GetUserRoles(userID)
		if err != nil {
			response.InternalError(c, "failed to load user roles")
			c.Abort()
			return
		}
		for _, role := range roles {
			if role.IsAdmin() {
				c.Next()
				return
			}
		}

		permissions, err := repo.GetUserPermissions(userID)
		if err != nil {
			response.InternalError(c, "failed to load user permissions")
			c.Abort()
			return
		}
		permissionMap := make(map[string]struct{}, len(permissions))
		for _, permission := range permissions {
			permissionMap[permission.Code] = struct{}{}
		}
		for _, code := range normalized {
			if _, ok := permissionMap[code]; ok {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "permission denied")
		c.Abort()
	}
}
