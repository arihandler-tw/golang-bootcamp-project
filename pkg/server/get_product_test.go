package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProduct(t *testing.T) {
	type args struct {
		url     string
		method  string
		inStore []map[string]any
	}

	tests := []struct {
		name      string
		args      args
		code      int
		isProduct bool
		res       map[string]any
	}{
		{
			name: "GET an existing product returns it",
			args: args{
				url:    "/product/id1",
				method: "GET",
				inStore: []map[string]any{
					{
						"id":          "id1",
						"price":       1.0,
						"description": "valid description",
					},
				},
			},
			code:      200,
			isProduct: true,
			res: map[string]any{
				"ID":          "id1",
				"Price":       1.0,
				"Description": "valid description",
			},
		},
		{
			name: "GET an inexistent products returns not found",
			args: args{
				url:    "/product/id1",
				method: "GET",
			},
			code:      404,
			isProduct: false,
			res: map[string]any{
				"error": "not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := SetupRoutes()

			populateStore(tt.args.inStore, router)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.args.method, tt.args.url, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			var actual map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &actual)
			if err != nil {
				t.Errorf("Unexpected error during unmarshaling of the response: %v", err)
			}
			if tt.isProduct {
				assert.Equal(t, true, compareProductResponses(actual, tt.res))
			} else {
				assert.Equal(t, actual, tt.res)
			}

		})
	}
}

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
