package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCorsMiddleware(t *testing.T) {
	setupRouter()
	router.Use(Cors())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "CORS test passed")
	})
	defer sqlDB.Close()

	server := httptest.NewServer(router)
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL+"/test", nil)
	req.Header.Set("Origin", "https://localhost")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Access-Control-Allow-Origin") != "https://localhost" {
		t.Errorf("Unexpected Access-Control-Allow-Origin: %s", resp.Header.Get("Access-Control-Allow-Origin"))
	}
}
