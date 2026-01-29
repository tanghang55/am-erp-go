package upload

import (
	"os"
	"path/filepath"
	"time"

	"am-erp-go/internal/infrastructure/response"

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
		response.BadRequest(c, "file required")
		return
	}

	subDir := c.PostForm("subDir")
	if subDir == "" {
		subDir = "products"
	}

	targetDir := filepath.Join(h.baseDir, subDir)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	filename := time.Now().Format("20060102150405") + "-" + filepath.Base(file.Filename)
	dst := filepath.Join(targetDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	url := "/" + filepath.ToSlash(filepath.Join(h.baseDir, subDir, filename))
	response.Success(c, map[string]interface{}{"url": url, "filename": filename})
}
