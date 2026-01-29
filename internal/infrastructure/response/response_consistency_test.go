package response

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandlersAvoidGinHResponses(t *testing.T) {
	root := repoRoot(t)

	var files []string
	files = append(files, listHandlerFiles(t, filepath.Join(root, "internal", "module"))...)
	files = append(files, listGoFiles(t, filepath.Join(root, "internal", "infrastructure", "upload"))...)

	var offenders []string
	for _, filePath := range files {
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read %s: %v", filePath, err)
		}

		if strings.Contains(string(content), "gin.H{") {
			offenders = append(offenders, filePath)
		}
	}

	if len(offenders) > 0 {
		t.Fatalf("found handlers using gin.H responses:\n%s", strings.Join(offenders, "\n"))
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working dir: %v", err)
	}

	return filepath.Clean(filepath.Join(wd, "..", "..", ".."))
}

func listHandlerFiles(t *testing.T, moduleRoot string) []string {
	t.Helper()

	entries, err := os.ReadDir(moduleRoot)
	if err != nil {
		t.Fatalf("failed to read module dir: %v", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		httpDir := filepath.Join(moduleRoot, entry.Name(), "delivery", "http")
		files = append(files, listGoFiles(t, httpDir)...)
	}

	return files
}

func listGoFiles(t *testing.T, dir string) []string {
	t.Helper()

	stat, err := os.Stat(dir)
	if err != nil || !stat.IsDir() {
		return nil
	}

	var files []string
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".go") && !strings.HasSuffix(d.Name(), "_test.go") {
			files = append(files, path)
		}

		return nil
	})

	return files
}
