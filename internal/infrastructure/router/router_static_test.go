package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/infrastructure/upload"

	"github.com/gin-gonic/gin"
)

func TestRouterServesUploads(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tempDir := t.TempDir()
	t.Setenv("UPLOAD_DIR", tempDir)
	filePath := filepath.Join(tempDir, "demo.txt")
	if err := os.WriteFile(filePath, []byte("ok"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	jwtManager := auth.NewJWTManager("test", 1)
	r := NewRouter(jwtManager, nil, nil, nil, nil, nil, nil, nil, nil, nil, upload.NewUploadHandler(), nil, nil)
	engine := r.Setup()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/uploads/demo.txt", nil)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "ok" {
		t.Fatalf("expected file content, got %q", w.Body.String())
	}
}
