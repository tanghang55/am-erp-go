package upload

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	baseDir string
}

func NewUploadHandler() *UploadHandler {
	baseDir := os.Getenv("UPLOAD_DIR")
	if baseDir == "" {
		baseDir = "uploads"
	}
	return &UploadHandler{baseDir: baseDir}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "file required"})
		return
	}

	subDir := c.PostForm("subDir")
	if subDir == "" {
		subDir = "products"
	}

	targetDir := filepath.Join(h.baseDir, subDir)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	filename := time.Now().Format("20060102150405") + "-" + filepath.Base(file.Filename)
	dst := filepath.Join(targetDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	url := "/" + filepath.ToSlash(filepath.Join(h.baseDir, subDir, filename))
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    gin.H{"url": url, "filename": filename},
	})
}
