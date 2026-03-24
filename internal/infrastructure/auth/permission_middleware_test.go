package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	identityDomain "am-erp-go/internal/module/identity/domain"

	"github.com/gin-gonic/gin"
)

type stubPermissionRepo struct {
	roles       []identityDomain.Role
	permissions []identityDomain.Permission
	err         error
}

func (s *stubPermissionRepo) GetUserRoles(userID uint64) ([]identityDomain.Role, error) {
	return s.roles, s.err
}

func (s *stubPermissionRepo) GetUserPermissions(userID uint64) ([]identityDomain.Permission, error) {
	return s.permissions, s.err
}

func TestRequirePermissionRejectsWhenPermissionMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &stubPermissionRepo{
		roles: []identityDomain.Role{
			{Name: "operator"},
		},
		permissions: []identityDomain.Permission{
			{Code: "inventory.view"},
		},
	}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(UserIDKey, uint64(8))
		c.Next()
	})
	router.GET("/secured", RequirePermission(repo, "sales.manage"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/secured", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestRequirePermissionAllowsAdminRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &stubPermissionRepo{
		roles: []identityDomain.Role{
			{Name: "admin"},
		},
	}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(UserIDKey, uint64(8))
		c.Next()
	})
	router.GET("/secured", RequirePermission(repo, "sales.manage"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/secured", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
