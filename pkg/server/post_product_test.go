package server

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostProduct(t *testing.T) {
	router := SetupRoutes()

	w := httptest.NewRecorder()
	payload, _ := json.Marshal(map[string]any{
		"price":       1.0,
		"description": "valid description",
	})
	req, _ := http.NewRequest(
		"POST",
		"/product/id1",
		bytes.NewReader(payload),
	)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	//assert.Equal(t, "pong", w.Body.String())
}
