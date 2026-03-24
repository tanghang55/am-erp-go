package upload

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
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
	var resp struct {
		Data struct {
			Filename string `json:"filename"`
			URL      string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	matched, err := regexp.MatchString(`^UPL\d{16}\.png$`, resp.Data.Filename)
	if err != nil {
		t.Fatalf("regexp error: %v", err)
	}
	if !matched {
		t.Fatalf("unexpected upload filename: %s", resp.Data.Filename)
	}
	if filepath.Base(resp.Data.URL) != resp.Data.Filename {
		t.Fatalf("expected url to end with filename, got url=%s filename=%s", resp.Data.URL, resp.Data.Filename)
	}
	if !regexp.MustCompile(`^/uploads/products/UPL\d{16}\.png$`).MatchString(resp.Data.URL) {
		t.Fatalf("unexpected upload url: %s", resp.Data.URL)
	}
}

func TestUploadImageUsesConfiguredUrlBase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tempDir := t.TempDir()
	t.Setenv("UPLOAD_DIR", tempDir)
	t.Setenv("UPLOAD_URL_BASE", "/static/uploads")

	handler := NewUploadHandler()
	router := gin.New()
	router.POST("/api/upload/image", handler.UploadImage)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("subDir", "products")
	part, _ := writer.CreateFormFile("file", "demo.jpg")
	part.Write([]byte("fake"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload/image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if !regexp.MustCompile(`^/static/uploads/products/UPL\d{16}\.jpg$`).MatchString(resp.Data.URL) {
		t.Fatalf("unexpected custom upload url: %s", resp.Data.URL)
	}
}

func TestResolveURLBase(t *testing.T) {
	t.Setenv("UPLOAD_URL_BASE", "static/uploads")
	if got := ResolveURLBase(); got != "/static/uploads" {
		t.Fatalf("expected /static/uploads, got %s", got)
	}
}
