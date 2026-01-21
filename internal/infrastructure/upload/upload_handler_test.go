package upload

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUploadImageReturnsUrl(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tempDir := t.TempDir()
	t.Setenv("UPLOAD_DIR", tempDir)

	handler := NewUploadHandler()
	router := gin.New()
	router.POST("/api/upload/image", handler.UploadImage)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("subDir", "products")
	part, _ := writer.CreateFormFile("file", "demo.png")
	part.Write([]byte("fake"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload/image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if _, err := os.Stat(filepath.Join(tempDir, "products")); err != nil {
		t.Fatalf("expected upload dir to exist")
	}
}
