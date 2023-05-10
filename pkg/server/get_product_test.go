package server_test

import (
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProduct(t *testing.T) {
	router := SetupRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/product/id1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
