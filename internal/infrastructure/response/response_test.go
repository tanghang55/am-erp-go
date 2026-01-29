package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func performRequest(t *testing.T, handler func(*gin.Context)) (int, map[string]interface{}) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler(c)

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	return w.Code, body
}

func TestSuccessResponse(t *testing.T) {
	status, body := performRequest(t, func(c *gin.Context) {
		Success(c, gin.H{"foo": "bar"})
	})

	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", status)
	}

	if body["code"].(float64) != 0 {
		t.Fatalf("expected code 0, got %v", body["code"])
	}

	if body["message"].(string) != "success" {
		t.Fatalf("expected message success, got %v", body["message"])
	}

	data := body["data"].(map[string]interface{})
	if data["foo"].(string) != "bar" {
		t.Fatalf("expected data.foo to be bar, got %v", data["foo"])
	}
}

func TestSuccessWithMessage(t *testing.T) {
	status, body := performRequest(t, func(c *gin.Context) {
		SuccessWithMessage(c, "created", gin.H{"id": 10})
	})

	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", status)
	}

	if body["message"].(string) != "created" {
		t.Fatalf("expected message created, got %v", body["message"])
	}
}

func TestPaginatedResponse(t *testing.T) {
	status, body := performRequest(t, func(c *gin.Context) {
		Paginated(c, []int{1, 2}, 2)
	})

	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", status)
	}

	data := body["data"].(map[string]interface{})
	if data["total"].(float64) != 2 {
		t.Fatalf("expected total 2, got %v", data["total"])
	}

	if _, ok := data["page"]; ok {
		t.Fatalf("did not expect page in Paginated response")
	}
}

func TestSuccessPageResponse(t *testing.T) {
	status, body := performRequest(t, func(c *gin.Context) {
		SuccessPage(c, []int{1, 2}, 2, 3, 50)
	})

	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", status)
	}

	data := body["data"].(map[string]interface{})
	if data["total"].(float64) != 2 {
		t.Fatalf("expected total 2, got %v", data["total"])
	}

	if data["page"].(float64) != 3 {
		t.Fatalf("expected page 3, got %v", data["page"])
	}

	if data["page_size"].(float64) != 50 {
		t.Fatalf("expected page_size 50, got %v", data["page_size"])
	}
}

func TestServerErrorResponse(t *testing.T) {
	status, body := performRequest(t, func(c *gin.Context) {
		ServerError(c, "boom")
	})

	if status != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", status)
	}

	if body["code"].(float64) != 500 {
		t.Fatalf("expected code 500, got %v", body["code"])
	}

	if body["message"].(string) != "boom" {
		t.Fatalf("expected message boom, got %v", body["message"])
	}
}
