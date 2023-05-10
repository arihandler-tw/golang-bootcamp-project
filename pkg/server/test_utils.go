package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func populateStore(store []map[string]any, router *gin.Engine) {
	w := httptest.NewRecorder()
	for _, p := range store {
		payload, _ := json.Marshal(map[string]any{
			"price":       p["price"],
			"description": p["description"],
		})
		url := fmt.Sprintf("%v%v", "/product/", p["id"])
		req, _ := http.NewRequest("POST", url, bytes.NewReader(payload))
		router.ServeHTTP(w, req)
	}
}

func compareProductResponses(p1, p2 map[string]any) bool {
	return p1["ID"] == p2["ID"] && p1["Price"] == p2["Price"] && p1["Description"] == p2["Description"]
}
