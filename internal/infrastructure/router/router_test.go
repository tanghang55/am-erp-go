package router

import (
	"testing"

	"am-erp-go/internal/infrastructure/upload"
)

func TestUploadStaticBaseUsesConfiguredURLBase(t *testing.T) {
	t.Setenv("UPLOAD_URL_BASE", "/static/uploads")
	if got := upload.ResolveURLBase(); got != "/static/uploads" {
		t.Fatalf("expected /static/uploads, got %s", got)
	}
}
