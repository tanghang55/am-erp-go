package upload

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"am-erp-go/internal/infrastructure/numbering"
	"am-erp-go/internal/infrastructure/response"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	baseDir string
	urlBase string
}

func ResolveURLBase() string {
	urlBase := os.Getenv("UPLOAD_URL_BASE")
	if urlBase == "" {
		urlBase = "/uploads"
	}
	return "/" + filepath.ToSlash(filepath.Clean(strings.TrimPrefix(urlBase, "/")))
}

func NewUploadHandler() *UploadHandler {
	baseDir := os.Getenv("UPLOAD_DIR")
	if baseDir == "" {
		baseDir = "uploads"
	}
	return &UploadHandler{
		baseDir: baseDir,
		urlBase: ResolveURLBase(),
	}
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

	ext := filepath.Ext(file.Filename)
	filename := numbering.Generate("UPL", time.Now()) + ext
	dst := filepath.Join(targetDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	url := h.urlBase + "/" + filepath.ToSlash(filepath.Join(subDir, filename))
	response.Success(c, map[string]interface{}{"url": url, "filename": filename})
}
