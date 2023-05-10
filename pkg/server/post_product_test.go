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
	type args struct {
		url     string
		method  string
		payload map[string]any
	}

	tests := []struct {
		name      string
		args      args
		code      int
		isProduct bool
		res       map[string]any
	}{
		{
			name: "Valid POST returns the product",
			args: args{
				url:    "/product/id1",
				method: "POST",
				payload: map[string]any{
					"price":       1.0,
					"description": "valid description",
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
			name: "Returns an error when product posted description is more than 50 char long",
			args: args{
				url:    "/product/id1",
				method: "POST",
				payload: map[string]any{
					"price":       1.0,
					"description": "invalid description because the description is over fifty characters",
				},
			},
			code:      500,
			isProduct: false,
			res: map[string]any{
				"error": "description should be less than 50 characters long",
			},
		},
		{
			name: "Returns an error when ID has non-ASCII characters",
			args: args{
				url:    "/product/id日本",
				method: "POST",
				payload: map[string]any{
					"price":       1.0,
					"description": "valid description",
				},
			},
			code:      500,
			isProduct: false,
			res: map[string]any{
				"error": "ID should contain only ASCII characters",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := SetupRoutes()

			w := httptest.NewRecorder()
			payload, _ := json.Marshal(tt.args.payload)
			req, _ := http.NewRequest(tt.args.method, tt.args.url, bytes.NewReader(payload))
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
