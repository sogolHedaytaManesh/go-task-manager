package http_test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	HTTPhandler "task-manager/internal/http"
	"testing"
)

func TestCORSMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(HTTPhandler.CORSMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))

	reqOptions, _ := http.NewRequest("OPTIONS", "/test", nil)
	wOptions := httptest.NewRecorder()
	router.ServeHTTP(wOptions, reqOptions)
	assert.Equal(t, http.StatusNoContent, wOptions.Code)
}

func TestTracingMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(HTTPhandler.TracingMiddleware())
	router.GET("/trace", func(c *gin.Context) {
		traceID := HTTPhandler.TraceIDFromContext(c.Request.Context())
		c.String(http.StatusOK, traceID)
	})

	req, _ := http.NewRequest("GET", "/trace", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	traceID := w.Body.String()
	assert.NotEmpty(t, traceID)
	assert.True(t, strings.HasPrefix(traceID, ""))
}
