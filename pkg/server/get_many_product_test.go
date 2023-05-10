package server

import (
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetManyProduct(t *testing.T) {
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
		res       []map[string]any
	}{
		{
			name: "GET 2 products returns the first two",
			args: args{
				url:    "/product?limit=2",
				method: "GET",
				inStore: []map[string]any{
					{
						"id":          "id1",
						"price":       1.0,
						"description": "valid description",
					},
					{
						"id":          "id2",
						"price":       1.0,
						"description": "valid description",
					},
					{
						"id":          "id3",
						"price":       1.0,
						"description": "valid description",
					},
				},
			},
			code:      200,
			isProduct: true,
			res: []map[string]any{
				{
					"ID":          "id1",
					"Price":       1.0,
					"Description": "valid description",
				},
				{
					"ID":          "id2",
					"Price":       1.0,
					"Description": "valid description",
				},
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
			var actual []map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &actual)
			if err != nil {
				t.Errorf("Unexpected error during unmarshaling of the response: %v", err)
			}
			if tt.isProduct {
				for i, r := range tt.res {
					assert.Equal(t, true, compareProductResponses(actual[i], r))
				}
			} else {
				assert.Equal(t, actual, tt.res)
			}

		})
	}
}
